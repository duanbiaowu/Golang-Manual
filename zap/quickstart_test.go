package zap

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSugar(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger) // flushes buffer, if any

	url := "https://www.uber.com"

	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)

	sugar.Infof("Failed to fetch URL: %s", url)
}

func TestLogger(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	url := "https://www.uber.com"

	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

func TestNamespace(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	logger.Info("tracked some metrics",
		zap.Namespace("metrics"),
		zap.Int("counter", 1),
	)

	logger2 := logger.With(
		zap.Namespace("metrics"),
		zap.Int("counter", 2),
	)
	logger2.Info("tracked some metrics")
}

func TestUseConfig(t *testing.T) {
	rawJSON := `
{
	"level": "debug",
	"encoding": "json",
	"outputPaths": ["stdout", "server.log"],
	"errorOutputPaths": ["stderr"],
	"initialFields": {"name": "Go"},
	"encoderConfig": {
	  "messageKey": "message",
	  "levelKey": "level",
	  "levelEncoder": "lowercase"
	}
}
`

	var cfg zap.Config
	if err := json.Unmarshal([]byte(rawJSON), &cfg); err != nil {
		t.Fatal(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	logger.Info("server start work successfully!")
}

func output(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func TestAddCallerSkip(t *testing.T) {
	// 有时稍微封装了一下记录日志的方法，但是希望输出的文件名和行号是调用封装函数的位置。
	// 可以使用 zap.AddCallerSkip(skip int) 向上跳 1 层：
	logger, _ := zap.NewProduction(zap.AddCaller(), zap.AddCallerSkip(1))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)

	output("hello world")
}

func f1() {
	f2("hello world")
}

func f2(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

func TestAddStackTrace(t *testing.T) {
	logger, _ := zap.NewProduction(zap.AddStacktrace(zapcore.WarnLevel))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)

	f1()
}

func TestReplaceGlobals(t *testing.T) {
	// 为了方便使用，zap提供了两个全局的Logger，一个是*zap.Logger，可调用zap.L()获得；另一个是*zap.SugaredLogger，
	// 可调用zap.S()获得。需要注意的是，全局的 Logger 默认并不会记录日志

	// 输出
	// gt -run='TestReplaceGlobals'
	// === RUN   TestReplaceGlobals
	// {"level":"info","msg":"global Logger after"}
	// {"level":"info","msg":"global SugaredLogger after"}
	// 调用ReplaceGlobals之前记录的日志并没有输出
	zap.L().Info("global Logger before")
	zap.S().Info("global SugaredLogger before")

	logger := zap.NewExample()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)
	zap.L().Info("global Logger after")
	zap.S().Info("global SugaredLogger after")
}

func TestCommonFields(t *testing.T) {
	// 如果每条日志都要记录一些共用的字段，那么使用 zap.Fields(fs ...Field)创建的选项。
	// 例如在服务器日志中记录可能都需要记录 serverId 和 serverName
	logger, _ := zap.NewProduction(zap.Fields(
		zap.Int("serverId", 1024),
		zap.String("serverName", "user_orders"),
	))

	logger.Info("hello world")
	logger.Info("hello world")
	logger.Info("hello world")
}

func TestWithNativeLog(t *testing.T) {
	// 如果项目一开始使用的是标准日志库log，后面想转为zap。
	// 不必修改之前的文件，可以调用zap.NewStdLog(l *Logger) *log.Logger返回一个标准的log.Logger，
	// 内部实际上写入的还是 zap.Logger
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	std := zap.NewStdLog(logger)
	std.Print("standard logger wrapper")

	undo := zap.RedirectStdLog(logger)
	log.Print("redirected standard library")
	undo()

	log.Print("restored standard library")
}

// See the doc: https://geektutu.com/post/hpg-reflect.html

package performance

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// 配置默认从 json 文件中读取，如果环境变量中设置了某个配置项，则以环境变量中的配置为准。
// 配置项和环境变量对应的规则非常简单：将 json 字段的字母转为大写，将 - 转为下划线，并添加 CONFIG_ 前缀。

type RefConfig struct {
	Name    string `json:"server-name"` // CONFIG_SERVER_NAME
	IP      string `json:"server-ip"`   // CONFIG_SERVER_IP
	URL     string `json:"server-url"`  // CONFIG_SERVER_URL
	Timeout string `json:"timeout"`     // CONFIG_TIMEOUT
}

func readRefConfig() *RefConfig {
	config := RefConfig{}
	typ := reflect.TypeOf(config)
	value := reflect.Indirect(reflect.ValueOf(&config))

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if v, ok := f.Tag.Lookup("json"); ok {
			key := fmt.Sprintf("CONFIG_%s", strings.ReplaceAll(strings.ToUpper(v), "-", "_"))
			if env, exist := os.LookupEnv(key); exist {
				value.FieldByName(f.Name).Set(reflect.ValueOf(env))
			}
		}
	}

	return &config
}

package golang

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var db struct {
	Dns string
}

// 如果测试文件中包含该函数，那么生成的测试将调用 TestMain(m)，而不是直接运行测试。
// TestMain 运行在主 goroutine 中 , 可以在调用 m.Run 前后做任何设置和拆卸。
// 注意，在 TestMain 函数的最后，应该使用 m.Run 的返回值作为参数去调用 os.Exit

// 在调用 TestMain 时 , flag.Parse 并没有被调用。
// 所以，如果 TestMain 依赖于 command-line 标志（包括 testing 包的标志），
// 则应该显式地调用 flag.Parse。
// 注意，这里的依赖是指，若 TestMain 函数内需要用到 command-line 标志，
// 则必须显式地调用 flag.Parse，否则不需要，因为 m.Run 中调用 flag.Parse。
func TestMain(m *testing.M) {
	db.Dns = os.Getenv("DATABASE_DNS")
	if db.Dns == "" {
		db.Dns = "root:123456@tcp(localhost:3306)/?charset=utf8&parseTime=True&loc=Local"
	}

	flag.Parse()
	exitCode := m.Run()

	db.Dns = ""
	fmt.Println("db.Dns = ", db.Dns)

	// 退出
	os.Exit(exitCode)
}

func TestDatabase(t *testing.T) {
	fmt.Println(db.Dns)
}

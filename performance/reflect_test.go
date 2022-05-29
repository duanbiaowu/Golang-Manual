package performance

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRefConfig(t *testing.T) {
	_ = os.Setenv("CONFIG_SERVER_NAME", "web_server")
	_ = os.Setenv("CONFIG_SERVER_IP", "127.0.0.1")
	_ = os.Setenv("CONFIG_SERVER_URL", "localhost")

	c := readRefConfig()

	assert.Equal(t, c.Name, "web_server")
	assert.Equal(t, c.IP, "127.0.0.1")
	assert.Equal(t, c.URL, "localhost")
}

func BenchmarkManualNewRefConfig(b *testing.B) {
	var config *RefConfig
	for i := 0; i < b.N; i++ {
		config = new(RefConfig)
	}
	_ = config
}

func BenchmarkReflectNewRefConfig(b *testing.B) {
	var config *RefConfig
	typ := reflect.TypeOf(RefConfig{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config, _ = reflect.New(typ).Interface().(*RefConfig)
	}
	_ = config
}

func BenchmarkSetRefConfig(b *testing.B) {
	config := new(RefConfig)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Name = "name"
		config.IP = "ip"
		config.URL = "url"
		config.Timeout = "timeout"
	}
}

func BenchmarkReflectFieldSetRefConfig(b *testing.B) {
	typ := reflect.TypeOf(RefConfig{})
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.Field(0).SetString("name")
		ins.Field(1).SetString("ip")
		ins.Field(2).SetString("url")
		ins.Field(3).SetString("timeout")
	}
}

func BenchmarkReflectFieldByNameSetRefConfig(b *testing.B) {
	typ := reflect.TypeOf(RefConfig{})
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.FieldByName("Name").SetString("name")
		ins.FieldByName("IP").SetString("ip")
		ins.FieldByName("URL").SetString("url")
		ins.FieldByName("Timeout").SetString("timeout")
	}
}

// 避免直接调用 FieldByName
// 可以利用字典将 Name 和 Index 的映射缓存起来。避免每次反复查找，耗费大量的时间。
func BenchmarkReflectFieldByNameCacheSetRefConfig(b *testing.B) {
	typ := reflect.TypeOf(RefConfig{})
	cache := make(map[string]int)
	for i := 0; i < typ.NumField(); i++ {
		cache[typ.Field(i).Name] = i
	}
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.Field(cache["Name"]).SetString("name")
		ins.Field(cache["IP"]).SetString("ip")
		ins.Field(cache["URL"]).SetString("url")
		ins.Field(cache["Timeout"]).SetString("timeout")
	}
}

// 使用反射给每个字段赋值，相比直接赋值，性能劣化约 100 - 1000 倍
// 反射 FieldByName 的比 Field 劣化 10 倍
//BenchmarkManualNewRefConfig
//BenchmarkManualNewRefConfig-8                           28240077                39.54 ns/op
//BenchmarkReflectNewRefConfig
//BenchmarkReflectNewRefConfig-8                          22649200                53.54 ns/op
//BenchmarkSetRefConfig
//BenchmarkSetRefConfig-8                                 1000000000               0.2924 ns/op
//BenchmarkReflectFieldSetRefConfig
//BenchmarkReflectFieldSetRefConfig-8                     61444561                19.74 ns/op
//BenchmarkReflectFieldByNameSetRefConfig
//BenchmarkReflectFieldByNameSetRefConfig-8                3548463               327.4 ns/op
//BenchmarkReflectFieldByNameCacheSetRefConfig
//BenchmarkReflectFieldByNameCacheSetRefConfig-8          23657766                50.20 ns/op

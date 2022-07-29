## go build
通过-tags命令行参数同时指定多个build标志，它们之间用空格分隔
```shell
# 只有在”linux/386“或”darwin平台下非cgo环境“才进行构建
# linux,386中linux和386用逗号链接表示AND的意思；而linux,386和darwin,!cgo之间通过空白分割来表示OR的意思
// +build linux,386 darwin,!cgo
```

## CGO中常用的功能
https://github.com/chai2010/cgo
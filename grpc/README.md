Protobuf核心的工具集是C++语言开发的，在官方的protoc编译器中并不支持Go语言。要想基于上面的hello.proto文件生成相应的Go代码，需要安装相应的插件。首先是安装官方的protoc工具，可以从 https://github.com/google/protobuf/releases 下载。然后是安装针对Go语言的代码生成插件，可以通过go get github.com/golang/protobuf/protoc-gen-go命令安装。

```shell
protoc --go_out=plugins=grpc:. hello.proto
```

## reference
1. https://chai2010.cn/advanced-go-programming-book/ch4-rpc/ch4-02-pb-intro.html
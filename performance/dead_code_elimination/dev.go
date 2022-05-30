//go:build dev
// +build dev

package main

const dev = true

// 条件编译
// //+build dev 表示 build tags 中包含 dev 时，该源文件参与编译。

// go build -o server -tags dev .
// ./server
// 2022/05/30 22:06:08 dev mode is enabled

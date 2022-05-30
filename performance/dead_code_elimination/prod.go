//go:build !dev
// +build !dev

package main

const dev = false

// 条件编译
// //+build !dev 表示 build tags 中不包含 dev 时，该源文件参与编译。

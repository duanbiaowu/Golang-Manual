package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// 1. Simple Example ----------------------------------------
func decorator(fn func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("Started")
		fn(s)
		fmt.Println("Done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}

// 调用了一个高阶函数
// 在调用时，先将 Hello() 函数传进去，然后返回一个匿名函数
// 匿名函数除了运行自己的代码，也调用了被传入 Hello() 函数
//decorator(Hello)("Hello World")

// 这样写可以增加代码可读性
//hello := decorator(Hello)
//hello("Hello World")

// 2. Simple Example ----------------------------------------
type SumFunc func(int64, int64) int64

func getFuncName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func timedSumFunc(f SumFunc) SumFunc {
	return func(start int64, end int64) int64 {
		// 这个参数设定很简洁
		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n", getFuncName(f), time.Since(t))
		}(time.Now())

		return f(start, end)
	}
}

func Sum1(start, end int64) int64 {
	var sum int64 = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i < end; i++ {
		sum += i
	}
	return sum
}

// 高斯算法
func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return (end - start + 1) * (end + start) / 2
}

//sum1 := timedSumFunc(Sum1)
//sum2 := timedSumFunc(Sum2)
//fmt.Printf("%d, %d\n", sum1(-10000, 10000000), sum2(-10000, 10000000))

// 3. HTTP Server ----------------------------------------
func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithServerHeader()")
		w.Header().Set("Server", "HelloServer v0.0.1")
		h(w, r)
	}
}

func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithAuthCookie()")
		cookie := &http.Cookie{Name: "Auth", Value: "Pass", Path: "/"}
		http.SetCookie(w, cookie)
		h(w, r)
	}
}
func WithBasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithBasicAuth()")
		cookie, err := r.Cookie("Auth")
		if err != nil || cookie.Value != "Pass" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w, r)
	}
}
func WithDebugLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithDebugLog")
		_ = r.ParseForm()
		log.Println(r.Form)
		log.Println("path", r.URL.Path)
		log.Println("scheme", r.URL.Scheme)
		log.Println(r.Form["url_long"])
		for k, v := range r.Form {
			log.Println("key:", k)
			log.Println("val:", strings.Join(v, ""))
		}
		h(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received Request %s from %s\n", r.URL.Path, r.RemoteAddr)
	_, _ = fmt.Fprintf(w, "hello world"+r.URL.Path)
}

//http.HandleFunc("/v1/hello", WithServerHeader(WithAuthCookie(hello)))
//http.HandleFunc("/v2/hello", WithServerHeader(WithBasicAuth(hello)))
//http.HandleFunc("/v3/hello", WithServerHeader(WithBasicAuth(WithDebugLog(hello))))
//err := http.ListenAndServe(":8080", nil)
//if err != nil {
//log.Fatal("ListenAndServe: ", err)
//}

// 3. HTTP Server ----------------------------------------
// 多个修饰器的 Pipeline
type HttpHandleDecorator func(http.HandlerFunc) http.HandlerFunc

func Handler(h http.HandlerFunc, decors ...HttpHandleDecorator) http.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-i-1] // iterate in reverse
		h = d(h)
	}

	return h
}
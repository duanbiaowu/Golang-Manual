# 注意事项（避坑指南）
## 1. 一个阻塞，全员等待
使用 singleflight 我们比较常见的是直接使用 Do 方法，但是这个极端情况下会导致整个程序 hang 住，如果我们的代码出点问题，有一个调用 hang 住了，
那么会导致所有的请求都 hang 住

```go
func singleFlightGetNumber(sg *singleflight.Group) int {
	v, _, _ := sg.Do("getNumber", func() (interface{}, error) {
        select {} // mock hanging bug
		return getNumber(), nil
	})
	return v.(int)
}
```
可以使用 DoChan 结合 select 做超时控制
```go
func singleFlightGetNumber(sg *singleflight.Group) int {
    v := sg.DoChan("getNumber", func() (interface{}, error) {
        select {} // mock hanging bug
        return getNumber(), nil
    })
	
    select {
    case r := <-v:
        return r.Val.(int)
    case <-time.After(time.Second * 3): // 也可以传入一个含 超时的 context 即可，执行时就会返回超时错误
        return 0
    }
}
```

## 2. 一个出错，全部出错
这个本身不是什么问题，因为 singleflight 就是这么设计的，但是实际使用的时候 如果我们一次调用要 1s，我们的数据库请求或者下游服务可以支撑 10rps 
的请求的时候这会导致我们的错误阈提高，因为实际上我们可以一秒内尝试 10 次，但是用了 singleflight 之后只能尝试一次，只要出错这段时间内的所有请求
都会受影响。这种情况我们可以启动一个 Goroutine 定时 forget 一下，相当于将 rps 从 1rps 提高到了 10rps


```go
go func() {
   time.Sleep(100 * time.Millisecond)
   // logging
   g.Forget(key)
}()
```

# reference
1. https://lailin.xyz/post/go-training-week5-singleflight.html
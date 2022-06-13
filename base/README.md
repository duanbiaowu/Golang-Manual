## goroutine
1. 请将是否异步调用的选择权交给调用者 
   1. 不然很有可能大家并不知道你在这个函数里面使用了 goroutine
2. 如果你要启动一个 goroutine 请对它负责
   1. 永远不要启动一个你无法控制它退出，或者你无法知道它何时推出的 goroutine
   2. 启动 goroutine 时请加上 panic recovery 机制，避免服务直接不可用
   3. 造成 goroutine 泄漏的主要原因就是 goroutine 中造成了阻塞，并且没有外部手段控制它退出
3. 尽量避免在请求中直接启动 goroutine 来处理问题
   1. 应该通过启动 worker 来进行消费，这样可以避免由于请求量过大，而导致大量创建 goroutine 从而导致 OOM

## reference
1. https://lailin.xyz/post/go-training-week3-goroutine.html
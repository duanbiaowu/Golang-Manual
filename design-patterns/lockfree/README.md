# 总结
这里的无锁其实只是没用互斥锁，用了原子操作，atomic 源码实际上在 CPU 上还是有锁的。

虽然在一些情况下 atomic 的性能要好很多，但这是一个 low level 的库，在实际的业务代码中最好还是使用 channel 。
但是我们也需要知道，在一些基础库，或者是需要极致性能的地方用上这个还是很爽的，但是使用的过程中一定要小心，不然还是会容易出 bug。


# reference
1. https://github.com/golang-design/lockfree
2. https://lailin.xyz/post/go-training-week3-atomic.html
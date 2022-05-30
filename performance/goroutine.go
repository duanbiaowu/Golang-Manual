// See the doc: https://geektutu.com/post/hpg-timeout-goroutine.html

package performance

import (
	"fmt"
	"time"
)

func doBadThing(done chan bool) {
	time.Sleep(time.Second)
	done <- true
}

func doGoodThing(done chan bool) {
	time.Sleep(time.Second)
	select {
	case done <- true:
	default:
		return
	}
}

// done 是一个无缓冲区的 channel，
// 如果没有超时，doBadThing 中会向 done 发送信号，select 中会接收 done 的信号，
// 	因此 doBadThing 能够正常退出，子协程也能够正常退出
// 当超时发生时，select 接收到 time.After 的超时信号就返回了，
// 	done 没有了接收方(receiver)，而 doBadThing 在执行 1s 后向 done 发送信号，
//	由于没有接收者且无缓存区，发送者(sender)会一直阻塞，导致协程不能退出
func timeout(f func(chan bool)) error {
	done := make(chan bool)
	go f(done)
	select {
	case <-done:
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):
		return fmt.Errorf("timeout")
	}
}

// 创建channel done 时，缓冲区设置为 1，即使没有接收方，发送方也不会发生阻塞。
func timeoutWithBuffer(f func(chan bool)) error {
	done := make(chan bool, 1)
	go f(done)
	select {
	case <-done:
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):
		return fmt.Errorf("timeout")
	}
}

// 分阶段复杂任务
func do2phases(phase1, done chan bool) {
	time.Sleep(time.Second) // 第 1 段
	select {
	case phase1 <- true:
	default:
		return
	}
	time.Sleep(time.Second) // 第 2 段
	done <- true
}

// 这里只能使用 select，而不能使用设置缓冲区的方式。
// 因为如果给 channel phase1 设置了缓冲区，phase1 <- true 总能执行成功，
//	那么无论是否超时，都会执行到第二阶段，而没有即时返回，这是我们不愿意看到的。
//	对应到上面的业务，就可能发生一种异常情况，向客户端发送了 2 次响应:
// 		<1 任务超时执行，向客户端返回超时，
//		<2 一段时间后，向客户端返回执行结果。
// 缓冲区不能够区分是否超时了，但是 select 可以（没有接收方，channel 发送信号失败，则说明超时了）。
func timeoutFirstPhase() error {
	phase1 := make(chan bool)
	done := make(chan bool)
	go do2phases(phase1, done)
	select {
	case <-phase1:
		<-done
		fmt.Println("done")
		return nil
	case <-time.After(time.Millisecond):
		return fmt.Errorf("timeout")
	}
}

// 强制 kill goroutine 可能吗？
// 答案是不能，goroutine 只能自己退出，而不能被其他 goroutine 强制关闭或杀死。
// goroutine 被设计为不可以从外部无条件地结束掉，只能通过 channel 来与它通信。也就是说，每一个 goroutine 都需要承担自己退出的责任。
//	(A goroutine cannot be programmatically killed. It can only commit a cooperative suicide.)

// 一些观点：
// 	<1 杀死一个 goroutine 设计上会有很多挑战，当前所拥有的资源如何处理？堆栈如何处理？defer 语句需要执行么？
// 	<2 如果允许 defer 语句执行，那么 defer 语句可能阻塞 goroutine 退出，这种情况下怎么办呢？

// 一些建议：
// 	<1 尽量使用非阻塞 I/O（非阻塞 I/O 常用来实现高性能的网络库），阻塞 I/O 很可能导致 goroutine 在某个调用一直等待，而无法正确结束。
//	<2 业务逻辑总是考虑退出机制，避免死循环。
//	<3 任务分段执行，超时后即时退出，避免 goroutine 无用的执行过多，浪费资源。

# 支持类型
* CPU ：CPU 分析，采样消耗 cpu 的调用，这个一般用来定位排查程序里耗费计算资源的地方；
* Memory ：内存分析，一般用来排查内存占用，内存泄露等问题；
* Block ：阻塞分析，会采样程序里阻塞的调用情况；
* Mutex ：互斥锁分析，采样互斥锁的竞争情况；

# 开启方式
* net/http/pprof ：使用在 web 服务器的场景；(只是在 runtime/pprof 上的一层 web 封装)
* runtime/pprof  ：使用在非服务器应用程序的场景；

# Memory

## 基础点
* golang 内存 pprof 是采样的，每 512KB 采样一次；
* golang 的内存采样的是堆栈路径，而不是类型信息；
* golang 的内存采样入口一定是通过mProf_Malloc，mProf_Free 这两个函数。所以，如果是 cgo 分配的内存，那么是没有机会调用到这两个函数的，所以如果是 cgo 导致的内存问题，go tool pprof 是分析不出来的；

## 分析形式
1. 如果是 net/http/pporf 方式开启的，那么可以直接在控制台上输入，浏览器就能看；
2. 另一种方式是先把信息 dump 到本地文件，然后用 go tool 去分析（生产环境通用）

## 指标含义备忘录
```shell
// https://github.com/golang/go/blob/master/src/runtime/mstats.go#L150

// 总共从OS申请的字节数
// 是下面各种XxxSys指标的总和。包含运行时的heap、stack和其他内部数据结构的总和。
// 它是虚拟内存空间。不一定全部映射成了物理内存。
Sys

// 见`Sys`
HeapSys

// 还在使用的对象，以及不使用还没被GC释放的对象的字节数
// 平时应该平缓，gc时可能出现锯齿
HeapAlloc

// 正在使用的对象字节数。
// 有个细节是，如果一个span中可包含多个object，只要一个object在使用，那么算的是整个span。
// `HeapInuse` - `HeapAlloc`是GC中保留，可以快速被使用的内存。
HeapInuse

// 已归还给OS的内存。没被堆再次申请的内存。
HeapReleased

// 没被使用的span的字节数。
// 这部分内存可以被归还给OS，并且还包含了`HeapReleased`。
// 可以被再次申请，甚至作为栈内存使用。
// `HeapIdle` - `HeapReleased`即GC保留的。
HeapIdle

/// ---

// 和`HeapAlloc`一样
Alloc

// 累计的`Alloc`
// 累计的意思是随程序启动后一直累加增长，永远不会下降。
TotalAlloc

// 没什么卵用
Lookups = 0

// 累计分配的堆对象数
Mallocs

// 累计释放的堆对象数
Frees

// 存活的对象数。见`HeapAlloc`
// HeapObjects = `Mallocs` - `Frees`
HeapObjects

// ---
// 下面的XxxInuse中的Inuse的含义，和XxxSys中的Sys的含义，基本和`HeapInuse`和`HeapSys`是一样的
// 没有XxxIdle，是因为都包含在`HeapIdle`里了

// StackSys基本就等于StackInuse，再加上系统线程级别的栈内存
Stack = StackInuse / StackSys

// 为MSpan结构体使用的内存
MSpan = MSpanInuse / MSpanSys

// 为MCache结构体使用的内存
MCache = MCacheInuse / MCacheSys

// 下面几个都是底层内部数据结构用到的XxxSys的内存统计
BuckHashSys
GCSys
OtherSys

// ---
// 下面是跟GC相关的

// 下次GC的触发阈值，当HeapAlloc达到这个值就要GC了
NextGC

// 最近一次GC的unix时间戳
LastGC

// 每个周期中GC的开始unix时间戳和结束unix时间戳
// 一个周期可能有0次GC，也可能有多次GC，如果是多次，只记录最后一个
PauseNs
PauseEnd

// GC次数
NumGC

// 应用程序强制GC的次数
NumForcedGC

// GC总共占用的CPU资源。在0~1之间
GCCPUFraction

// 没被使用，忽略就好 
DebugGC

```

**查看方式**
```go
// 方式一
import "runtime"

var m runtime.MemStats
runtime.ReadMemStats(&m)
    
// 方式二
import _ "net/http/pprof"
import "net/http"
    
http.ListenAndServe("0.0.0.0:10001", nil)
// http://127.0.0.1:10001/debug/pprof/heap?debug=1
// curl -sS 'http://127.0.0.1:10001/debug/pprof/heap?seconds=5' -o heap.pporf
```

## 示例
```shell
# 运行测试
go test -v -run='TestPprof' .
# dump 采样
curl -sS 'http://127.0.0.1:10001/debug/pprof/heap?seconds=5' -o heap.pporf
```

## 结果说明
```shell
go tool pprof ./29075_20190523_154406_heap
(pprof) o              
...          
  sample_index              = inuse_space          //: [alloc_objects | alloc_space | inuse_objects | inuse_space]
...       
(pprof) alloc_space
(pprof) top
Showing nodes accounting for 290MB, 100% of 290MB total
      flat  flat%   sum%        cum   cum%
     140MB 48.28% 48.28%      140MB 48.28%  main.funcA (inline)
     100MB 34.48% 82.76%      190MB 65.52%  main.funcB (inline)
      50MB 17.24%   100%      140MB 48.28%  main.funcC (inline)
         0     0%   100%      290MB   100%  main.main
         0     0%   100%      290MB   100%  runtime.main
```

**这个 top 信息表明了这么几点信息：**
* main.funcA  这个函数现场分配了 140M 的内存，main.funcB 这个函数现场分配了 100M 内存，main.funcC 现场分配了 50M 内存；
  * 现场的意思：纯粹自己函数直接分配的，而不是调用别的函数分配的； 
  * 这些信息通过 flat 得知；
* main.funcA  分配的 140M 内存纯粹是自己分配的，没有调用别的函数分配过内存；
  * 这个信息通过 main.funcA flat 和 cum 都为 140 M 得出；
* main.funcB  自己分配了 100MB，并且还调用了别的函数，别的函数里面涉及了 90M 的内存分配；
  * 这个信息通过 main.funcB flat 和 cum 分别为 100 M，190M 得出；
* main.funcC  自己分配了 50MB，并且还调用了别的函数，别的函数里面涉及了 90M 的内存分配；
  * 这个信息通过 main.funcC flat 和 cum 分别为 50 M，140 M 得出；
* main.main ：所有分配内存的函数调用都是走这个函数出去的。main 函数本身没有函数分配，但是他调用的函数分配了 290M；


# reference
1. https://pengrl.com/p/20031/
2. https://mp.weixin.qq.com/s/OXpWRiCHcxpFlylUk3RDEA
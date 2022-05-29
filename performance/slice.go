// See the doc: https://geektutu.com/post/hpg-slice.html

package performance

import (
	"math/rand"
	"time"
)

// 在已有切片的基础上进行切片，不会创建新的底层数组。
// 因为原来的底层数组没有发生变化，内存会一直占用，直到没有变量引用该数组。
// 因此很可能出现一种极端情况，原切片由大量的元素构成，
// 但是我们在原切片的基础上切片，虽然只使用了很小一段，但底层数组在内存中仍然占据了大量空间，得不到释放。
// 比较推荐的做法，使用 copy 替代 re-slice。

func lastNumsBySlice(origin []int) []int {
	return origin[len(origin)-2:]
}

func lastNumsByCopy(origin []int) []int {
	result := make([]int, 2)
	copy(result, origin[len(origin)-2:])
	return result
}

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

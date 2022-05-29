// See the doc: https://geektutu.com/post/hpg-range.html

package performance

import (
	"math/rand"
	"time"
)

type RangeItem struct {
	id  int
	val [4096]byte
}

func rangeInternalGenerateInt(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func rangeInternalGenerateItem(n int) []*RangeItem {
	items := make([]*RangeItem, n)
	for i := 0; i < n; i++ {
		items[i] = &RangeItem{id: i}
	}
	return items
}

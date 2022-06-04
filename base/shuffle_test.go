// See the doc: https://chai2010.cn/advanced-go-programming-book/ch6-cloud/ch6-05-load-balance.html

package base

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func shuffle1(nums []int) {
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
}

func shuffle2(nums []int) {
	for i := len(nums); i > 0; i-- {
		lastIdx := i - 1
		idx := rand.Intn(i)
		nums[lastIdx], nums[idx] = nums[idx], nums[lastIdx]
	}
}

func TestShuffle(t *testing.T) {
	const N = 1000000

	var cnt1 = map[int]int{}
	for i := 0; i < N; i++ {
		var sl = []int{0, 1, 2, 3, 4, 5, 6}
		shuffle1(sl)
		cnt1[sl[0]]++
	}
	for _, v := range cnt1 {
		assert.Equal(t, "0.14", fmt.Sprintf("%.2f", float64(v)/float64(N)))
	}

	var cnt2 = map[int]int{}
	for i := 0; i < N; i++ {
		var sl = []int{0, 1, 2, 3, 4, 5, 6}
		shuffle2(sl)
		cnt2[sl[0]]++
	}

	for _, v := range cnt2 {
		assert.Equal(t, "0.14", fmt.Sprintf("%.2f", float64(v)/float64(N)))
	}
}

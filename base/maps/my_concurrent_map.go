package maps

import "sync"

const PRIME32 = uint32(16777619)

type ConcurrentMap struct {
	ShardCount int
	items      []*ConcurrentMapShard
}

type ConcurrentMapShard struct {
	sync.RWMutex
	m map[string]interface{}
}

func New(log2OfShardCount int) ConcurrentMap {
	shardCount := 1 << log2OfShardCount
	m := ConcurrentMap{
		ShardCount: shardCount,
		items:      make([]*ConcurrentMapShard, shardCount),
	}
	for i := 0; i < shardCount; i++ {
		m.items[i] = &ConcurrentMapShard{m: make(map[string]interface{})}
	}
	return m
}

func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShard {
	// 性能优化方式：n mod m = n & (m - 1)
	return m.items[uint(fnv32(key))&uint(m.ShardCount-1)]
}

func (m ConcurrentMap) Set(key string, value interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.m[key] = value
	shard.Unlock()
}

func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.m[key]
	shard.RUnlock()
	return val, ok
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	length := len(key)
	for i := 0; i < length; i++ {
		hash *= PRIME32
		hash ^= uint32(key[i])
	}
	return hash
}

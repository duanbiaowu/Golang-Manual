package cache

import (
	"testing"

	"github.com/golang/groupcache"
	"github.com/stretchr/testify/assert"
)

type gcMockDB struct {
	data map[string]string
}

func (db *gcMockDB) Get(key string) string {
	return db.data[key]
}

func (db *gcMockDB) Set(key string, value string) {
	db.data[key] = value
}

func NewGcMockDB() *gcMockDB {
	ndb := new(gcMockDB)
	ndb.data = make(map[string]string)
	return ndb
}

func TestGroupCacheGet(t *testing.T) {
	db := NewGcMockDB()

	db.Set("foo", "bar")
	db.Set("one", "two")

	var stringGroup = groupcache.NewGroup("SlowDBCache", 10<<1, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result := db.Get(key)
			err := dest.SetString(result)
			if err != nil {
				return err
			}
			return nil
		}))

	var err error
	var data []byte

	err = stringGroup.Get(nil, "foo", groupcache.AllocatingByteSliceSink(&data))
	assert.Nil(t, err)
	err = stringGroup.Get(nil, "one", groupcache.AllocatingByteSliceSink(&data))
	assert.Nil(t, err)

	db.Set("foo", "bar2")
	err = stringGroup.Get(nil, "foo", groupcache.AllocatingByteSliceSink(&data))
	assert.Nil(t, err)

	assert.Equal(t, "bar", string(data))
}

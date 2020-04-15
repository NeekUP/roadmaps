package infrastructure

import (
	"github.com/NeekUP/roadmaps/core"
	"github.com/patrickmn/go-cache"
	"time"
)

type inmemoryCache struct {
	memcache *cache.Cache
}

func NewInMemoryCache() core.DistributedCache {
	return &inmemoryCache{memcache: cache.New(10*time.Minute, 20*time.Minute)}
}

func (mem *inmemoryCache) Save(key string, item interface{}, duration time.Duration) error {
	return mem.memcache.Add(key, item, duration)
}

func (mem *inmemoryCache) Get(key string) (interface{}, bool) {
	item, ok := mem.memcache.Get(key)
	return item, ok
}

func (mem *inmemoryCache) Delete(key string) {
	mem.memcache.Delete(key)
}

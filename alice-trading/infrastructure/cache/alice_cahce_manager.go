package cache

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

// キャッシュ管理（Singleton）
type AliceCacheManager struct {
	AliceCache *cache.Cache
}

var cacheManager *AliceCacheManager
var once sync.Once

// CacheManagerを取得します。
func GetCacheManager() AliceCacheManager {
	once.Do(func() {
		cacheManager = &AliceCacheManager{}
		var manager = AliceCacheManager{AliceCache: cache.New(30*time.Minute, 60*time.Minute)}
		*cacheManager = manager
	})
	return *cacheManager
}

// Keyを指定してキャッシュを保存します。
func (a *AliceCacheManager) Set(key string, target interface{}, d enum.Duration) {
	switch d {
	case enum.DEFAULT:
		a.AliceCache.Set(key, target, cache.DefaultExpiration)
		break
	case enum.NONE:
		a.AliceCache.Set(key, target, cache.NoExpiration)
		break
	default:
		panic("something wrong duration")
	}
}

// Keyに一致するキャッシュを取得します。
func (a *AliceCacheManager) Get(key string) interface{} {
	target, ok := a.AliceCache.Get(key)
	if !ok {
		log.Print("no cache for", key)
	}
	return target
}

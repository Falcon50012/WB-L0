package my_cache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	OrderCache *cache.Cache
	CacheMx    sync.Mutex
)

func InitCache() {
	OrderCache = cache.New(5*time.Minute, 10*time.Minute)
}

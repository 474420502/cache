package cache

import (
	"log"
	"sync"
	"time"
)

// CacheBlock 缓存
type CacheBlock struct {
	share interface{}

	updateMehtod  UpdateMehtod
	onUpdateError func(err interface{})

	lastUpdate time.Time
	interval   time.Duration

	vLock sync.Mutex
	value interface{}
}

// NewBlockCache 创建一个Cache对象
func NewBlockCache(interval time.Duration, u UpdateMehtod) Cache {

	c := &CacheBlock{
		updateMehtod: u,
		interval:     interval,
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}

	// c.update()
	return c
}

func (cache *CacheBlock) SetShare(share interface{}) {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	cache.share = share
}

// SetOnUpdateError 默认false
func (cache *CacheBlock) SetOnUpdateError(errFunc func(err interface{})) {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	cache.onUpdateError = errFunc
}

// Destroy 异步更新必须调用Destroy, 销毁对象
func (cache *CacheBlock) Destroy() {

}

// update 主动更新 没锁
func (cache *CacheBlock) update() {

	defer func() {
		cache.lastUpdate = time.Now()
		if err := recover(); err != nil {
			cache.onUpdateError(err)
		}
	}()

	func() {
		value := cache.updateMehtod(cache.share)
		if value == nil {
			return
		}

		if err, ok := value.(error); ok {
			cache.onUpdateError(err)
			return
		}
		cache.value = value
	}()
}

// Value 获取缓存的值
func (cache *CacheBlock) Value() interface{} {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()

	if time.Since(cache.lastUpdate) >= cache.interval {
		cache.update()
	}

	return cache.value
}

// Value 获取缓存的值
func (cache *CacheBlock) GetUpdate() time.Time {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	return cache.lastUpdate
}

package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

// CacheInterval 缓存
type CacheInterval struct {
	slock sync.Mutex
	share interface{}

	isDestroy int32

	updateMehtod UpdateMehtod

	onUpdateError func(err interface{})

	lastUpdate time.Time
	interval   time.Duration

	vLock sync.Mutex
	value interface{}
}

func (cache *CacheInterval) SetShare(share interface{}) {
	cache.slock.Lock()
	defer cache.slock.Unlock()
	cache.share = share
}

// SetOnUpdateError 默认false
func (cache *CacheInterval) SetOnUpdateError(errFunc func(err interface{})) {
	cache.slock.Lock()
	defer cache.slock.Unlock()
	cache.onUpdateError = errFunc
}

// Destroy 异步更新必须调用Destroy, 销毁对象
func (cache *CacheInterval) Destroy() {
	atomic.StoreInt32(&cache.isDestroy, 1)
}

// update 主动更新 没锁
func (cache *CacheInterval) update() {
	cache.slock.Lock()
	defer cache.slock.Unlock()

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

		cache.vLock.Lock()
		cache.value = value
		cache.vLock.Unlock()
	}()
}

// Value 获取缓存的值
func (cache *CacheInterval) Value() interface{} {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	return cache.value
}

// Value 获取缓存的值
func (cache *CacheInterval) GetUpdate() time.Time {
	cache.slock.Lock()
	defer cache.slock.Unlock()
	return cache.lastUpdate
}

package cache

import (
	"log"
	"sync"
	"time"
)

// CacheBlock 缓存 对比 默认 CacheInterval 的. 该方法有阻塞效果. 就是更新过程会阻塞
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
func NewBlockCache(interval time.Duration, updateMethod UpdateMehtod) Cache {
	c := &CacheBlock{
		updateMehtod: updateMethod,
		interval:     interval,
		onUpdateError: func(err interface{}) {
			log.Println(err)
		},
	}
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

// update 主动更新 没锁全阻塞
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

// ForceUpdate 强制更新
func (cache *CacheBlock) ForceUpdate() {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	cache.update()
}

// Value 获取缓存的值
func (cache *CacheBlock) GetUpdate() time.Time {
	cache.vLock.Lock()
	defer cache.vLock.Unlock()
	return cache.lastUpdate
}

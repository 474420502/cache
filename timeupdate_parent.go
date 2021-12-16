package cache

import "time"

// CacheInterval 加一层为了 自动析构 自动释放线程
type CacheInterval struct {
	cache *cacheIntervalBackup
}

func (cache *CacheInterval) SetShare(share interface{}) {
	cache.cache.SetShare(share)
}

// SetOnUpdateError 默认false
func (cache *CacheInterval) SetOnUpdateError(errFunc func(err interface{})) {
	cache.cache.SetOnUpdateError(errFunc)
}

// Value 获取缓存的值
func (cache *CacheInterval) Value() interface{} {

	return cache.cache.Value()
}

// Value 获取缓存的值
func (cache *CacheInterval) ForceUpdate() {
	cache.cache.ForceUpdate()
}

// GetUpdate 获取缓存的值
func (cache *CacheInterval) GetUpdate() time.Time {
	return cache.cache.GetUpdate()
}

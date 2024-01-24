package cache

// // cacheIntervalBackup 缓存
// type cacheIntervalBackup struct {
// 	slock sync.Mutex
// 	share interface{}

// 	isDestroy     chan struct{} // 是否退出持续更新
// 	isForceUpdate chan struct{} // 是否强制更新 无视时间间隔

// 	updateMehtod UpdateMehtod

// 	onUpdateError func(err interface{})

// 	lastUpdate time.Time
// 	interval   *time.Ticker

// 	vLock sync.Mutex
// 	value interface{}
// }

// func (cache *cacheIntervalBackup) SetShare(share interface{}) {
// 	cache.slock.Lock()
// 	defer cache.slock.Unlock()
// 	cache.share = share
// }

// // SetOnUpdateError 默认false
// func (cache *cacheIntervalBackup) SetOnUpdateError(errFunc func(err interface{})) {
// 	cache.slock.Lock()
// 	defer cache.slock.Unlock()
// 	cache.onUpdateError = errFunc
// }

// // Destroy 异步更新必须调用Destroy, 销毁对象
// func (cache *cacheIntervalBackup) Destroy() {
// 	cache.isDestroy <- struct{}{}
// }

// // update 主动更新 没锁
// func (cache *cacheIntervalBackup) update() {
// 	cache.slock.Lock()
// 	defer cache.slock.Unlock()

// 	defer func() {
// 		cache.lastUpdate = time.Now()
// 		if err := recover(); err != nil {
// 			cache.onUpdateError(err)
// 		}
// 	}()

// 	func() {
// 		value := cache.updateMehtod(cache.share)
// 		if value == nil {
// 			return
// 		}
// 		if err, ok := value.(error); ok {
// 			cache.onUpdateError(err)
// 			return
// 		}

// 		cache.vLock.Lock()
// 		cache.value = value
// 		cache.vLock.Unlock()
// 	}()
// }

// // ValueSync 同步获取缓存的值, 其他操作必须等待完成
// func (cache *cacheIntervalBackup) ValueSync(do func(v interface{})) {
// 	cache.vLock.Lock()
// 	defer cache.vLock.Unlock()
// 	do(cache.value)
// }

// // Value 获取缓存的值
// func (cache *cacheIntervalBackup) Value() interface{} {
// 	cache.vLock.Lock()
// 	defer cache.vLock.Unlock()
// 	return cache.value
// }

// // Value 获取缓存的值
// func (cache *cacheIntervalBackup) ForceUpdate() {
// 	cache.isForceUpdate <- struct{}{}
// 	<-cache.isForceUpdate
// }

// // Value 获取缓存的值
// func (cache *cacheIntervalBackup) GetUpdate() time.Time {
// 	cache.slock.Lock()
// 	defer cache.slock.Unlock()
// 	return cache.lastUpdate
// }

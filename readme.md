## Cache 
用于定时或者条件更新数据. 例如: 1.爬虫请求页面时间或者某个值 2.请求某表的内容信息, 减少频繁操作

## 使用方法

```golang
cache := New(func() interface{} {
    resp, _ := gcurl.Execute(`curl "http://httpbin.org/uuid"`)
    return string(resp.Content())
}) // 创建 cache 并添加请求方法. 默认为异步更新. 当到触发条件的时间. 返回值还是上次值. 更新完后再次调用才是最新值. cache.SetBlock 可以用这个函数设置为阻塞. 更新会让更新方法彻底完成后. 返回值.
cache.SetUpdateInterval(time.Millisecond * 50) // 设置更新时间
// cache.Update() 主动更新 
log.Println(cache.Value()) // { "uuid": "fbbea7fc-3a71-4e65-bc18-6642eddf83a7" }
time.Sleep(time.Millisecond * 50)
log.Println(cache.Value()) // { "uuid": "125c0d4d-e598-4b37-b0da-c928d44b703d" }
```
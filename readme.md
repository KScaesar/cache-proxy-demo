# 傳遞 message 來撰寫程式, 在 golang 應該如何思考

實現一個 簡易的 cache proxy  
藉由  concurrency 的情景, 提供數種程式的撰寫方式  

![go meetup profile](./asset/go%20meetup%20profile.jpg)

## slide

v2:  
2023-04-25 go meetup 所使用的 slide  
<https://docs.google.com/presentation/d/1BKdpu8wF9zqpoGQrjQW6QmzOVIKwwwz6I96lvoMEPPg>

v0:  
內容分散, 關注太多細節, 程式碼凌亂  
~~<https://docs.google.com/presentation/d/1yctNKOoct49OEj7jZtKfVjrnZifABWbWxfpJ3MM2D9w>~~  

v0 相關程式  
<https://github.com/KScaesar/cache-proxy-demo/tree/v0.1.0>

## cache proxy implimentation

[base proxy](./cache_proxy.go)

- [global lock proxy](./mutex_proxy.go)
- [channel proxy](./channel_proxy.go)
- [shard lock proxy](./syncMap_proxy.go)
- [singleflight proxy](./singleflight_proxy.go)

![channel impl data flow](./asset/channel%202.gif)

## reference

1. [Go语言圣经中文版 - 9.7. 示例: 并发的非阻塞缓存](https://github.com/gopl-zh/gopl-zh.github.com/blob/master/ch9/ch9-07.md?fbclid=IwAR0sVeVwXrDVxT0Ozh0vcSTxVJV-scl_ZA-vCDFkJE9HqiyRBDkSrnOpWc8)
2. [Messaging Patterns - Return Address](https://www.enterpriseintegrationpatterns.com/patterns/messaging/ReturnAddress.html)
4. [sync.Map的LoadOrStore用途](https://xnum.github.io/2018/11/syncmap-loadorstore/)
6. [sync.singleflight 到底怎么用才对？](https://www.cyningsun.com/01-11-2021/golang-concurrency-singleflight.html)
8. [Hardware Memory Models](https://research.swtch.com/hwmm)
9. [How Does Golang Channel Works](https://levelup.gitconnected.com/how-does-golang-channel-works-6d66acd54753)
10. [Mutex Or Channel](https://github.com/golang/go/wiki/MutexOrChannel)

## benchmark

```bash
# OneKey data size = 2e4
BenchmarkMutexProxy-8            6037	200725 ns/op
BenchmarkChannelProxy-8          4477	260086 ns/op
BenchmarkSyncMapProxy-8         15058	 75422 ns/op
BenchmarkSingleflightProxy-8     4350	264185 ns/op

# OneKey data size = 4e4
BenchmarkMutexProxy-8            6176	201015 ns/op
BenchmarkChannelProxy-8          4332	264898 ns/op
BenchmarkSyncMapProxy-8         15966	 73374 ns/op
BenchmarkSingleflightProxy-8     4310	259505 ns/op

# MultiKey data size = 2e4
BenchmarkMutexProxy-8            4516	332631 ns/op
BenchmarkChannelProxy-8         62859	 24878 ns/op
BenchmarkSyncMapProxy-8         62914	 25624 ns/op
BenchmarkSingleflightProxy-8    65200	 19896 ns/op

# MultiKey data size = 4e4
BenchmarkMutexProxy-8            4928	307945 ns/op
BenchmarkChannelProxy-8         38280	 33544 ns/op
BenchmarkSyncMapProxy-8         31879	 44351 ns/op
BenchmarkSingleflightProxy-8    34802	 40129 ns/op
```

## 關於我

有問題討論, 可發 issue, 或用下方的聯絡方式

Email:  
x246libra@hotmail.com

Telegram id:  
@ksCaesar
# 傳遞 message 來撰寫程式, 在 golang 應該如何思考

實現一個 簡易的 cache proxy  
藉由  concurrency 的情景, 提供數種程式的撰寫方式  

## slide

v2: 2023-04-25 go meetup 所使用的 side  
<>

v0: 內容分散, 關注太多細節, 程式碼凌亂  
<https://docs.google.com/presentation/d/1yctNKOoct49OEj7jZtKfVjrnZifABWbWxfpJ3MM2D9w>  

v0 相關程式  
<https://github.com/KScaesar/cache-proxy-demo/tree/v0.1.0>

## cache proxy implimentation

- [global lock proxy](./mutex_proxy.go)
- [channel proxy](./channel_proxy.go)
- [shard lock proxy](./syncMap_proxy.go)
- [singleflight proxy](./singleflight_proxy.go)

![channel impl data flow](./asset/channel%202.gif)

## reference

1. [Go语言圣经中文版 - 9.7. 示例: 并发的非阻塞缓存](https://github.com/gopl-zh/gopl-zh.github.com/blob/master/ch9/ch9-07.md?fbclid=IwAR0sVeVwXrDVxT0Ozh0vcSTxVJV-scl_ZA-vCDFkJE9HqiyRBDkSrnOpWc8)
2. [Messaging Patterns - Return Address](https://www.enterpriseintegrationpatterns.com/patterns/messaging/ReturnAddress.html)
3. [Messaging Patterns Overview](https://www.enterpriseintegrationpatterns.com/patterns/messaging/)
4. [sync.Map的LoadOrStore用途](https://xnum.github.io/2018/11/syncmap-loadorstore/)
5. [golang-pkg - Singleflight and its usage in 17 Media](https://github.com/golangtw/GolangTaiwanGathering/blob/master/meetup/gtg51/slides/singleflight-for-meetup.pdf)
6. [sync.singleflight 到底怎么用才对？](https://www.cyningsun.com/01-11-2021/golang-concurrency-singleflight.html)
7. [clean-arch (Tung 東東)](https://docs.google.com/presentation/d/1ouNiohGRcl5m_uGNrwlHuZ_hAXH13joLGTtkkxyJ8eY/edit#slide=id.g1c2a9713f29_0_1)
8. [Hardware Memory Models](https://research.swtch.com/hwmm)
9. [How Does Golang Channel Works](https://levelup.gitconnected.com/how-does-golang-channel-works-6d66acd54753)
10. [Mutex Or Channel](https://github.com/golang/go/wiki/MutexOrChannel)

## benchmark

```bash

```

## 關於我

有問題討論, 可發 issue, 或用下方的聯絡方式

Email:  
x246libra@hotmail.com

Telegram id:  
@ksCaesar
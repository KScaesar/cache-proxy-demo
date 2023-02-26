# 傳遞 message 來撰寫程式, 在 golang 是怎樣的體驗

實現一個 簡易的 cache proxy  
藉由  concurrency 的情景, 提供數種程式的撰寫方式  
並探討單一職責  

## slide

<https://docs.google.com/presentation/d/1yctNKOoct49OEj7jZtKfVjrnZifABWbWxfpJ3MM2D9w>

## reference

- [Go语言圣经中文版 - 9.7. 示例: 并发的非阻塞缓存](https://github.com/gopl-zh/gopl-zh.github.com/blob/master/ch9/ch9-07.md?fbclid=IwAR0sVeVwXrDVxT0Ozh0vcSTxVJV-scl_ZA-vCDFkJE9HqiyRBDkSrnOpWc8)
- [Messaging Patterns - Return Address](https://www.enterpriseintegrationpatterns.com/patterns/messaging/ReturnAddress.html)
- [Messaging Patterns Overview](https://www.enterpriseintegrationpatterns.com/patterns/messaging/)
- [sync.Map的LoadOrStore用途](https://xnum.github.io/2018/11/syncmap-loadorstore/)
- [golang-pkg - Singleflight and its usage in 17 Media](https://github.com/golangtw/GolangTaiwanGathering/blob/master/meetup/gtg51/slides/singleflight-for-meetup.pdf)
- [sync.singleflight 到底怎么用才对？](https://www.cyningsun.com/01-11-2021/golang-concurrency-singleflight.html)
- [clean-arch (Tung 東東)](https://docs.google.com/presentation/d/1ouNiohGRcl5m_uGNrwlHuZ_hAXH13joLGTtkkxyJ8eY/edit#slide=id.g1c2a9713f29_0_1)

## benchmark

```bash

```

## 關於我

有問題討論, 可發 issue, 或用下方的聯絡方式

Email:  
x246libra@hotmail.com

Telegram id:  
@ksCaesar
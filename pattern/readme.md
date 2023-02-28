# Template 設計模式練習

由於 go 沒有繼承  
要實現 design pattern 中的 Template 模式  
跟 oop 語言的寫法會有所不同  
參考以下網站的寫法, 感到很奇怪  
<https://refactoring.guru/design-patterns/template-method/go/example>

打算在這個目錄, 利用 cache proxy 寫一個  
自己認為符合 go風格的 Template 設計模式  
和 root 目錄的寫法, 可能會有所不同  

主要概念是  
不要把 Template 執著以 物件 method 的方式實現  
簡單用 function 表達  
但缺點是 需要建造很多小型物件, 來滿足 interface  

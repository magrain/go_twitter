package main

import (
  "fmt"
  "net/http"
)

type String string

// String に ServeHTTP 関数を追加
func (s String) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, s)
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, World")
}

func main() {
  http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる

  http.HandleFunc("/tl", func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprint(w, "TLみたーい!!")
  })


  // ディレクトリ名とURLパスを変える場合
  // 例：http://localhost:8080/mysecret/sample1.txt
  http.Handle("/mypage/", http.StripPrefix("/mypage/", http.FileServer(http.Dir("../contents"))))

  //http.Handle("/tw", String("こんちわ"))

  http.ListenAndServe(":8080", nil)
}

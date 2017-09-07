package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "fmt"
)

func apiClockHandler(w http.ResponseWriter, r *http.Request) {
    // JSONにする構造体
    type ResponseBody struct {
        Time time.Time `json:"time"`
    }
    rb := &ResponseBody{
        Time: time.Now(),
    }

    // ヘッダーをセット
    w.Header().Set("Content-type", "application/json")

    // JSONにエンコードしてレスポンスに書き込む
    if err := json.NewEncoder(w).Encode(rb); err != nil {
        log.Fatal(err)
    }
}

func main() {
    // ハンドラーを登録
    http.HandleFunc("/api/clock", apiClockHandler)

    http.Handle("/mypage/", http.StripPrefix("/mypage/", http.FileServer(http.Dir("./contents_test"))))
    // サーバーを起動
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "起動OK")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
/*    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
*/
}

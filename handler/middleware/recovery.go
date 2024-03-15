package middleware

import (
	"log"
	"net/http"
)

// パニックハンドラ
func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		// 回復不可能なエラーをハンドリングする
		// recover()を呼ぶと，巻き戻しを停止して,panic()を処理できる
		// stack unwind中かな
		// 巻き戻し中に実行できるのは，deferで遅延指定された関数内だけ
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

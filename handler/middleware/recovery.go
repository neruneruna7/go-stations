package middleware

import (
	"log"
	"net/http"
)

// パニックハンドラ
func Recovery(h http.Handler) http.Handler {

	// 今回のコードだと，これはNewRouterが呼ばれた時に実行される
	// クロージャ内はリクエストが来た時に動くことがわかるように，こういうロギングをつけてみた
	log.Println("Recovery Middlware started")

	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Closure Recovery Middlware started")

		// TODO: ここに実装をする
		// 回復不可能なエラーをハンドリングする
		// recover()を呼ぶと，巻き戻しを停止して,panic()を処理できる
		// stack unwind中かな
		// 巻き戻し中に実行できるのは，deferで遅延指定された関数内だけ
		defer func() {
			log.Println("deferd Recovery Middlware started")

			if err := recover(); err != nil {

				log.Println("Error: ", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			log.Println("deferd Recovery Middlware finished")

		}()

		h.ServeHTTP(w, r)
		log.Println("Closure Recovery Middlware finished")
	}
	log.Println("Recovery Middlware finished")
	return http.HandlerFunc(fn)
}

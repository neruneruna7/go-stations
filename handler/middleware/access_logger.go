package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LogContent struct {
	Timestamp time.Time `json:"timestamp"`
	// 単位ミリ秒の処理時間
	Latency int64  `json:"latency"`
	Path    string `json:"path"`
	Os      string `json:"os"`
}

func AccessLogger(h http.Handler) http.Handler {
	log.Println("AccessLogger Middlware started")

	// 今回のコードだと，これはNewRouterが呼ばれた時に実行される
	// クロージャ内はリクエストが来た時に動くことがわかるように，こういうロギングをつけてみた

	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Closure AccessLogger Middlware started")

		// アクセスしてきたときの時間を記録する
		var startTime = time.Now()

		defer func() {
			log.Println("deferd AccessLogger Middlware started")
			var endTime = time.Now()
			var latency = endTime.Sub(startTime).Milliseconds()
			var path = r.URL.Path
			var os = r.Context().Value(CTX_OS_KEY).(string)

			var logContent = &LogContent{
				Timestamp: startTime,
				Latency:   latency,
				Path:      path,
				Os:        os,
			}

			serialized, err := json.Marshal(logContent)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(serialized))

			log.Println("deferd AccessLogger Middlware finished")
		}()

		h.ServeHTTP(w, r)
		log.Println("Closure AccessLogger Middlware finished")
	}
	log.Println("AccessLogger Middlware finished")
	return http.HandlerFunc(fn)
}

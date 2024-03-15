package middleware

import (
	"net/http"
)

func Middlwares(h http.Handler) http.Handler {
	// いちいちミドルウェアをすべて適用するのが面倒だから，まとめて適用する関数を作る

	fn := func(w http.ResponseWriter, r *http.Request) {
		var middlewares_wraped_handler = Recovery(CaptureDeviceOs(h))

		middlewares_wraped_handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

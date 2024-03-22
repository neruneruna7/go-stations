package middleware

import (
	"log"
	"net/http"
)

func CommonMiddlwares(h http.Handler) http.Handler {
	log.Println("CommonMiddlwares Middlware started")
	// いちいちミドルウェアをすべて適用するのが面倒だから，まとめて適用する関数を作る

	// fn := func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("Closure CommonMiddlwares Middlware started")
	// 	// deferはLIFOで積みあがっていくことに注意（スタック）
	// 	// Recoverを最初に入れる（deferは最後に実行する）にすることで，他のミドルウェア内のpanic
	// 	// も拾うことができるはず
	// 	var middlewares_wraped_handler = Recovery(CaptureDeviceOs(AccessLogger(h)))

	// 	middlewares_wraped_handler.ServeHTTP(w, r)
	// 	log.Println("Closure CommonMiddlwares Middlware finished")
	// }
	// log.Println("CommonMiddlwares Middlware finished")

	var middlewaresWrapedHandler = Recovery(CaptureDeviceOs(AccessLogger(h)))
	// return http.HandlerFunc(fn)
	return middlewaresWrapedHandler
}

package middleware

import (
	"crypto/subtle"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

func BasicAuth(b *model.BasicAuthConfig, h http.Handler) http.Handler {
	log.Println("BasicAuth  Middlware started")

	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Closure BasicAuth  Middlware started")

		var userid, password, ok = r.BasicAuth()
		if !ok {
			// このokはなんなんだ
			// ソースをみたら，basicauthのヘッダーのところをパースできたか否かっぽい
			// パース失敗だから問答無用で認証失敗でいいな
			log.Println("BasicAuth failed: parse failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// タイミング攻撃対策
		var isEqUserId = subtle.ConstantTimeCompare([]byte(userid), []byte(b.BasicAuthUserId))
		var isEqPassword = subtle.ConstantTimeCompare([]byte(password), []byte(b.BasicAuthPassword))

		if isEqUserId != 1 || isEqPassword != 1 {
			log.Println("BasicAuth failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
		log.Println("Closure BasicAuth  Middlware finished")
	}
	log.Println("BasicAuth Middlware finished")
	return http.HandlerFunc(fn)
}

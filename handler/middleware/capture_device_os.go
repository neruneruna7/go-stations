package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mileusna/useragent"
)

// キーの衝突を避けるため，独自の型をパッケージごとに定義してつかうのが推奨のようだ
// 非公開型であること
type ctxKey string

const CTX_OS_KEY ctxKey = "OS"

// リクエストを送ってきたデバイスのOSを取得する
func CaptureDeviceOs(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var user_agent_str = r.UserAgent()
		var ua = useragent.Parse(user_agent_str)
		var ua_os = ua.OS
		log.Println("hey!")
		log.Println(ua_os)

		// ブロック内でシャドーイングできないのがもどかしい
		// 解決策はあるはずなのだが
		var r2 = SetOs(r, ua_os)

		{
			// ContextにOS名が登録されているか確認する
			val, ok := r.Context().Value(CTX_OS_KEY).(string)
			fmt.Println(val, ok)
		}

		// 連鎖する以上，引数でとっているハンドラ（h）を，何かしら処理あるいは次に渡す必要があるはず
		// それがこのServeHTTPかな？
		// であるならば，コード内でのこのserveHttpの実行場所には，時間軸の関係があるはず
		h.ServeHTTP(w, r2)
	}
	return http.HandlerFunc(fn)
}

func SetOs(r *http.Request, os string) *http.Request {
	// データの流れがわかりづらい...
	// メモリの所有権はどこにあるんだ？？
	// http.Requestのポインタを取ってるんだから，そのポインタを指すメモリを書き換えているのか
	// そうでないのかみにくい

	ctx := r.Context()
	// Contexはポインタじゃないしなぁ
	ctx = context.WithValue(ctx, CTX_OS_KEY, os)
	r = r.WithContext(ctx)

	// でも返しておけば確実か
	return r
}

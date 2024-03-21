package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("main: started")
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	basic_auth_config := basic_auth_config_loader()

	log.Println("Loading PORT env")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("Loading DB_PATH")
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	log.Println("Setting time zone")
	// set time zone
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}
	log.Println("Time zone: ", time.Local)

	log.Println("Setting up db")
	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB, basic_auth_config)

	log.Println("running port:", port)
	// TODO: サーバーをlistenする
	// portを束縛する方法があるはず
	server := http.Server{
		Addr:    port,
		Handler: mux,
	}

	// linixはカーネルからいろいろシグナル来ることもあるしね
	// どのシグナルとか割り込みに対応するかはあとで考えよう
	// とりあえず，Ctrl+Cへの対応
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	// このstopの役割がわからん
	// deferをつけずにやったら起動して即終了したな
	// signal receivedってログには来てるし，この関数の最後まで突き抜けたのかな
	defer stop()

	// ゴルーチンでサーバーを起動する
	go func() {
		server.ListenAndServe()
	}()

	log.Println("awaiting signal")
	<-ctx.Done()
	// <-はチャンネルへ値を送信するといういことらしい
	// ということはメインチャンネルに値を送信してる？
	// 受け取る変数がないということは，事実上待機してるのと同じかな
	// ctx.Done()の内部では，NotifyContexで設定したシグナルが来るまで休眠してると思われる
	log.Println("signal received")
	// 5秒で終了するようにタイムアウトを設定
	// でも，5秒待たずに終了することがある
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 別にctxはgoroutinに渡してないし，このcancelはmain groutinを終わらせるのか？

	// もしや，待つまでもなくシャットダウンが完了して，この関数が終了したとき
	// このdefer cancel()が待つのを終わらせるのか？
	// 実行されないようにしても待つ必要がなければすぐ終了したな
	// 必要なら待ってたし 役割はなんだろう
	// deferをつけずにやったらシグナル即終了になったな
	// goroutinのリソースを解放するものらしいが...
	defer cancel()
	// cancel()

	// defer stop()
	// defer cancel()
	// 両方動作しないようにすると，シグナルで待たずに即終了するようになった．
	// どちらか片方でも動作すると，シグナル待ちになった
	// いやまて．どっちもコメントアウトでもシグナル待ちになるわ
	// 多分このShutdownでlisten停止命令を送るのと，停止待ち，タイムアウトを過ぎたら強制終了をやってるはず
	err = server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func basic_auth_config_loader() *model.BasicAuthConfig {
	var basic_auth_user_id = os.Getenv("BASIC_AUTH_USER_ID")
	var basic_auth_password = os.Getenv("BASIC_AUTH_PASSWORD")
	var basic_auth_config = model.BasicAuthConfig{
		BasicAuthUserId:   basic_auth_user_id,
		BasicAuthPassword: basic_auth_password,
	}
	return &basic_auth_config
}

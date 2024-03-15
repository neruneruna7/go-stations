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

	var basic_auth_user_id = os.Getenv("BASIC_AUTH_USER_ID")
	var basic_auth_password = os.Getenv("BASIC_AUTH_PASSWORD")
	var basic_auth_config = model.BasicAuthConfig{
		BasicAuthUserId:   basic_auth_user_id,
		BasicAuthPassword: basic_auth_password,
	}

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
	mux := router.NewRouter(todoDB, &basic_auth_config)

	log.Println("running port:", port)
	// TODO: サーバーをlistenする
	// portを束縛する方法があるはず
	var server = http.Server{
		Addr:    port,
		Handler: mux,
	}

	// linixはカーネルからいろいろシグナル来ることもあるしね
	// どのシグナルとか割り込みに対応するかはあとで考えよう
	// とりあえず，Ctrl+Cへの対応
	var ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		server.ListenAndServe()
	}()

	log.Println("awaiting signal")
	<-ctx.Done()
	log.Println("signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

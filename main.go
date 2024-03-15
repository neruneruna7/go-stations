package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
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
	var err error
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
	mux := router.NewRouter(todoDB)

	log.Println("running port:", port)
	// TODO: サーバーをlistenする
	// portを束縛する方法があるはず
	// NewRouterのドキュメントが見当たらない...?
	// TechTrain独自のやつか？
	http.ListenAndServe(port, mux)

	return nil
}

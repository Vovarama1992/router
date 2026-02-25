package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"router/internal/domain"
	"router/internal/infra"
	"router/internal/telegram"
)

func main() {
	ctx := context.Background()

	// ---- DB ----
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// ---- infra ----
	peerRepo := infra.NewPeerRepo(db)

	// ---- domain ----
	vpnService := domain.NewService(peerRepo)

	// ---- telegram ----
	app, err := telegram.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	bot := telegram.NewBot(app, vpnService)
	go app.Run(ctx, bot.Handle)

	addr := ":8080"
	log.Println("router: listening on", addr)

	select {}

}

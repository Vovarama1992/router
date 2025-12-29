package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"router/internal/delivery"
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

	// ---- http ----
	vpnHandler := delivery.NewVPNHandler(vpnService)

	r := chi.NewRouter()
	delivery.RegisterVPNRoutes(r, vpnHandler)

	addr := ":8080"
	log.Println("router: listening on", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"router/internal/config"
	"router/internal/delivery"
	"router/internal/domain"
	"router/internal/infra"
)

func main() {
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
	cfg := config.Load()
	peerRepo := infra.NewPeerRepo(db)

	vpnService := domain.NewService(cfg, peerRepo)

	// ---- handlers ----
	vpnHandler := delivery.NewVPNHandler(vpnService)

	// ---- router ----
	r := chi.NewRouter()
	delivery.RegisterVPNRoutes(r, vpnHandler)

	addr := ":8080"
	log.Println("router: listening on", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

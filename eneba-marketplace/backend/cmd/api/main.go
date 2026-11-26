package main

import (
	"context"
	"eneba-marketplace/backend/internal/config"
	"eneba-marketplace/backend/internal/database"
	"eneba-marketplace/backend/internal/server"
	"log"
	"net/http"
	"time"
)

func main() {
	c := config.Load()
	if c.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required; copy .env.example to .env")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, e := database.Open(ctx, c.DatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	s := http.Server{Addr: ":" + c.Port, Handler: server.Router(db, c.FrontendURL), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("API listening on port %s", c.Port)
	log.Fatal(s.ListenAndServe())
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/peterblog/blog/internal/api"
	"github.com/peterblog/blog/internal/auth"
	"github.com/peterblog/blog/internal/db"
)

func main() {
	// Load .env.local in development (no-op if file doesn't exist)
	_ = godotenv.Load(".env.local")

	auth.InitSessionStore()

	ctx := context.Background()
	if err := auth.InitProviders(ctx); err != nil {
		log.Printf("Warning: failed to init OIDC providers: %v", err)
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	router := api.NewRouter(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kaanchinar/url-shortener/handler"
	"github.com/kaanchinar/url-shortener/repo"
	"github.com/kaanchinar/url-shortener/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	urlRepo := repo.NewURLRepository(pool)
	urlService := service.NewURLService(urlRepo)
	baseURL := os.Getenv("BASE_URL")
	urlHandler := handler.NewURLHandler(urlService, baseURL)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/", func(r chi.Router) {
		r.With(httprate.LimitByIP(30, 1*time.Minute)).Post("/shorten", urlHandler.ShortenURL)
		r.Get("/{id}", urlHandler.GetLongURL)
	})

	log.Printf("Server running on %s\n", baseURL)
	log.Fatal(http.ListenAndServe(":3000", r))
}

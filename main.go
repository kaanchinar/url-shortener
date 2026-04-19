package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kaanchinar/url-shortener/handler"
	"github.com/kaanchinar/url-shortener/repo"
	"github.com/kaanchinar/url-shortener/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	urlRepo := repo.NewURLRepository(pool)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("pong"))
			if err != nil {
				return
			}
		})
		r.Post("/shorten", urlHandler.ShortenURL)
		r.Get("/{id}", urlHandler.GetLongURL)
	})

	fmt.Println("Server running on http://localhost:3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}

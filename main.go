package main

import (
	"log"
	"net/http"
	"os"

	"github.com/deathstarset/books-api/routes"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load(".env")
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Logger)
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mainRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health is good"))
	})
	mainRouter.Mount("/api/v1", routes.UsersRouter(client))
	mainRouter.Mount("/api/v1/books", routes.BooksRouter(client))

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port string not found")
	}

	server := &http.Server{Handler: mainRouter, Addr: ":" + portString}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}
}

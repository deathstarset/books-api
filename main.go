package main

import (
	"net/http"

	"github.com/deathstarset/books-api/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {

	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.Logger)
	mainRouter.Mount("/api/v1", routes.UsersRouter())

	http.ListenAndServe(":3000", mainRouter)
}

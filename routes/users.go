package routes

import (
	"github.com/deathstarset/books-api/handlers"
	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func UsersRouter(client *redis.Client, queries *database.Queries) *chi.Mux {
	usersRouter := chi.NewRouter()
	usersRouter.Post("/signup", middlewares.QueriesMiddleware(handlers.CreateUserHandler, queries))
	usersRouter.Post("/login", handlers.LoginHandler(client, queries))
	usersRouter.Get("/logout", handlers.LogoutHandler(client))
	return usersRouter
}

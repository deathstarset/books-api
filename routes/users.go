package routes

import (
	"github.com/deathstarset/books-api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func UsersRouter(client *redis.Client) *chi.Mux {
	usersRouter := chi.NewRouter()
	usersRouter.Post("/signup", apiConfig.CreateUserHandler)
	usersRouter.Post("/login", apiConfig.LoginHandler(client))
	usersRouter.Get("/logout", apiConfig.LogoutHandler(client))
	usersRouter.Get("/test", middlewares.AuthMiddleware(apiConfig.CreateBookHandler, client))
	return usersRouter
}

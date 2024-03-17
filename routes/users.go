package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func UsersRouter(client *redis.Client) *chi.Mux {
	usersRouter := chi.NewRouter()
	usersRouter.Post("/signup", apiConfig.CreateUserHandler)
	usersRouter.Post("/login", apiConfig.LoginHandler(client))
	usersRouter.Get("/logout", apiConfig.LogoutHandler(client))
	return usersRouter
}

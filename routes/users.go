package routes

import "github.com/go-chi/chi/v5"

func UsersRouter() *chi.Mux {
	usersRouter := chi.NewRouter()
	return usersRouter
}

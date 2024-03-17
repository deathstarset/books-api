package routes

import (
	"github.com/deathstarset/books-api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func BooksRouter(client *redis.Client) *chi.Mux {
	booksRouter := chi.NewRouter()
	booksRouter.Post("/", middlewares.AuthMiddleware(apiConfig.CreateBookHandler, client))
	booksRouter.Delete("/{bookID}", middlewares.AuthMiddleware(apiConfig.DeleteBookHandler, client))
	booksRouter.Get("/", apiConfig.GetAllBooksHandler)
	return booksRouter
}

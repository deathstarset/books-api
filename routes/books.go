package routes

import (
	"github.com/deathstarset/books-api/handlers"
	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func BooksRouter(client *redis.Client, queries *database.Queries) *chi.Mux {
	booksRouter := chi.NewRouter()
	booksRouter.Post("/", middlewares.AuthMiddleware(handlers.CreateBookHandler, client, queries))
	booksRouter.Delete("/{bookID}", middlewares.AuthMiddleware(handlers.DeleteBookHandler, client, queries))
	booksRouter.Patch("/{bookID}", middlewares.AuthMiddleware(handlers.UpdateBookHandler, client, queries))
	booksRouter.Get("/", middlewares.QueriesMiddleware(handlers.GetAllBooksHandler, queries))
	return booksRouter
}

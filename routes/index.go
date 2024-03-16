package routes

import (
	"github.com/deathstarset/books-api/handlers"
	"github.com/deathstarset/books-api/internal/database"
	"github.com/deathstarset/books-api/sql"
)

var apiConfig = handlers.ApiConfig{
	Queries: database.New(sql.ConnectDb()),
}

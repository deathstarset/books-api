package sql

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ConnectDb() *sql.DB {
	godotenv.Load(".env")
	connectionString := os.Getenv("DB_STRING")
	if connectionString == "" {
		log.Fatal("Database url not found")
	}

	connection, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	return connection
}

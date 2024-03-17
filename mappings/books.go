package mappings

import (
	"time"

	"github.com/deathstarset/books-api/internal/database"
	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageLink   string    `json:"image_link"`
	UserID      uuid.UUID `json:"user_id"`
}

func DbBookToBookMapping(book database.Book) Book {
	return Book{
		ID:          book.ID,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
		Title:       book.Title,
		Description: book.Description,
		ImageLink:   book.ImageLink,
		UserID:      book.UserID,
	}
}

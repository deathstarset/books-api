// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: books.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createBook = `-- name: CreateBook :one
INSERT INTO books (id, created_at, updated_at, title, description, image_link, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, title, description, image_link, user_id
`

type CreateBookParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description string
	ImageLink   string
	UserID      uuid.UUID
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	row := q.db.QueryRowContext(ctx, createBook,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Description,
		arg.ImageLink,
		arg.UserID,
	)
	var i Book
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
		&i.ImageLink,
		&i.UserID,
	)
	return i, err
}
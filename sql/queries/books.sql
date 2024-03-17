-- name: CreateBook :one
INSERT INTO books (id, created_at, updated_at, title, description, image_link, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
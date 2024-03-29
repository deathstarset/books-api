-- name: CreateBook :one
INSERT INTO books (id, created_at, updated_at, title, description, image_link, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books 
WHERE id = $1 AND user_id = $2;

-- name: GetAllBooks :many
SELECT * FROM books 
ORDER BY created_at DESC;

-- name: GetBookByUserAndID :one
SELECT * FROM books
WHERE id = $1 AND user_id = $2;

-- name: UpdateBook :exec
UPDATE books 
SET updated_at = $1, title = $2, description = $3
WHERE id = $4 AND user_id = $5;
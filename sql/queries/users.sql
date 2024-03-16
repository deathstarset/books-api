-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, username, password, pfp_link)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
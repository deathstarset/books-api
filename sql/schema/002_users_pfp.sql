-- +goose Up
ALTER TABLE users
ADD COLUMN pfp_link VARCHAR(255) UNIQUE NOT NULL
DEFAULT 'image-url.jpeg';

-- +goose Down
ALTER TABLE users DROP COLUMN pfp_link;
-- +goose Up
CREATE TABLE users(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  email VARCHAR NOT NULL UNIQUE CHECK(email ~*'^[A-Za-z0-9._]+@[A-Za-z0-9]+\\.[A-Za-z]{2,6}'),
  username VARCHAR NOT NULL UNIQUE CHECK(LENGTH(username) >= 2),
  password VARCHAR NOT NULL CHECK(LENGTH(password) >= 8)
);

-- +goose Down
DROP TABLE users;
-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, username, hashed_password, first_name, last_name)
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?;


-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

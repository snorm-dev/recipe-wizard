-- name: CreateIngredient :exec
INSERT INTO ingredients(id, created_at, updated_at, name, description)
VALUES (?, ?, ?, ?, ?);

-- name: GetIngredient :one
SELECT * FROM ingredients
WHERE id = ?;
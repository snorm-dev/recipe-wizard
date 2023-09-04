-- name: CreateRecipe :exec
INSERT INTO recipes(id, created_at, updated_at, name, description, url, prep_time, cook_time, total_time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetRecipe :one
SELECT * FROM recipes
WHERE id = ?;

-- name: GetRecipes :many
SELECT * FROM recipes;
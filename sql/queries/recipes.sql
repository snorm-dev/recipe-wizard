-- name: CreateRecipe :execresult
INSERT INTO recipes(created_at, updated_at, name, description, url, prep_time, cook_time, total_time, owner_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetRecipe :one
SELECT * FROM recipes
WHERE id = ?;

-- name: GetRecipesForUser :many
SELECT * FROM recipes
WHERE owner_id = ?;

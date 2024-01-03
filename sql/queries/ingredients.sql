-- name: CreateIngredient :execresult
INSERT INTO ingredients(created_at, updated_at, name, description, amount, units, standard_amount, standard_units, recipe_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetIngredient :one
SELECT * FROM ingredients
WHERE id = ?;

-- name: GetIngredientsForRecipe :many
SELECT * FROM ingredients
WHERE recipe_id = ?
ORDER BY id;

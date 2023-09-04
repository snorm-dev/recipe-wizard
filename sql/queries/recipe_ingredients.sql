-- name: CreateRecipeIngredient :exec
INSERT INTO recipe_ingredients(id, created_at, updated_at, ingredient_id, recipe_id, quantity, units)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetRecipeIngredient :one
SELECT * FROM recipe_ingredients
WHERE id = ?;

-- name: GetIngredientsForRecipe :many
SELECT sqlc.embed(ri), sqlc.embed(i) FROM recipe_ingredients ri
LEFT JOIN ingredients i ON i.id = ri.ingredient_id
WHERE ri.recipe_id = ?;
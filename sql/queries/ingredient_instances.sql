-- name: CreateIngredientInstance :execresult
INSERT INTO ingredient_instances (created_at, updated_at, ingredient_id, grocery_list_id, recipe_instance_id)
VALUES (?, ?, ?, ?, ?);

-- name: GetIngredientInstance :one
SELECT * FROM ingredient_instances
WHERE id = ?;

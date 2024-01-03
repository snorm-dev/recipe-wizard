-- name: CreateRecipeInstance :execresult
INSERT INTO recipe_instances (created_at, updated_at, grocery_list_id, recipe_id)
VALUES (?, ?, ?, ?);

-- name: GetRecipeInstance :one
SELECT * FROM recipe_instances
WHERE id = ?;

-- name: GetRecipeInstancesInGroceryList :many
SELECT * FROM recipe_instances ri 
WHERE ri.grocery_list_id = ?;

-- name: GetIngredientInstancesForRecipeInstance :many
SELECT * FROM ingredient_instances ii 
WHERE ii.recipe_instance_id = ?;

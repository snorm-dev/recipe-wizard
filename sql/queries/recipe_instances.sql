-- name: AddRecipeToGroceryList :execresult
INSERT INTO recipe_instances (created_at, updated_at, grocery_list_id, recipe_id)
VALUES (?, ?, ?, ?);

-- name: GetRecipeInstance :one
SELECT * FROM recipe_instances
WHERE id = ?;

-- name: GetRecipesInGroceryList :many
SELECT r.* FROM recipe_instances ri 
JOIN recipes r ON r.id = ri.recipe_id
WHERE ri.grocery_list_id = ?;

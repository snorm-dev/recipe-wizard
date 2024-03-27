-- name: CreateRecipeInstance :one
INSERT INTO recipe_instances (created_at, updated_at, grocery_list_id, recipe_id)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetRecipeInstance :one
SELECT * FROM recipe_instances
WHERE id = ?;

-- name: GetRecipeInstancesInGroceryList :many
SELECT * FROM recipe_instances ri 
WHERE ri.grocery_list_id = ?;

-- name: GetExtendedRecipeInstance :one
SELECT sqlc.embed(ri), sqlc.embed(r) from recipe_instances ri 
JOIN recipes r ON ri.recipe_id = r.id
WHERE ri.id = ?;

-- name: GetExtendedRecipeInstancesInGroceryList :many
SELECT sqlc.embed(ri), sqlc.embed(r) from recipe_instances ri 
JOIN recipes r ON ri.recipe_id = r.id
WHERE ri.grocery_list_id = ?;

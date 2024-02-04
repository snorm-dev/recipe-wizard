-- name: CreateItem :execresult
INSERT INTO ingredient_instances (created_at, updated_at, ingredient_id, grocery_list_id, recipe_instance_id)
VALUES (?, ?, ?, ?, ?);

-- name: GetItem :one
SELECT * FROM ingredient_instances
WHERE id = ?;

-- name: GetItemsForRecipeInstance :many
SELECT * FROM ingredient_instances ii 
WHERE ii.recipe_instance_id = ?;

-- name: GetExtendedItem :one
SELECT sqlc.embed(ii), sqlc.embed(i) FROM ingredient_instances ii
JOIN ingredients i ON ii.ingredient_id = i.id
WHERE ii.id = ?;

-- name: GetExtendedItemsForRecipeInstance :many
SELECT sqlc.embed(ii), sqlc.embed(i) FROM ingredient_instances ii
JOIN ingredients i ON ii.ingredient_id = i.id
WHERE ii.recipe_instance_id = ?;

-- name: GetItemsForGroceryList :many
SELECT * FROM ingredient_instances ii 
WHERE ii.grocery_list_id = ?;

-- name: GetExtendedItemsForGroceryList :many
SELECT sqlc.embed(ii), sqlc.embed(i) FROM ingredient_instances ii
JOIN ingredients i ON ii.ingredient_id = i.id
WHERE ii.grocery_list_id = ?;

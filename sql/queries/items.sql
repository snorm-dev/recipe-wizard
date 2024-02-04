-- name: CreateItem :execresult
INSERT INTO items (created_at, updated_at, ingredient_id, grocery_list_id, recipe_instance_id)
VALUES (?, ?, ?, ?, ?);

-- name: GetItem :one
SELECT * FROM items
WHERE id = ?;

-- name: GetItemsForRecipeInstance :many
SELECT * FROM items it 
WHERE it.recipe_instance_id = ?;

-- name: GetExtendedItem :one
SELECT sqlc.embed(it), sqlc.embed(i) FROM items it
JOIN ingredients i ON it.ingredient_id = i.id
WHERE it.id = ?;

-- name: GetExtendedItemsForRecipeInstance :many
SELECT sqlc.embed(it), sqlc.embed(i) FROM items it
JOIN ingredients i ON it.ingredient_id = i.id
WHERE it.recipe_instance_id = ?;

-- name: GetItemsForGroceryList :many
SELECT * FROM items it 
WHERE it.grocery_list_id = ?;

-- name: GetExtendedItemsForGroceryList :many
SELECT sqlc.embed(it), sqlc.embed(i) FROM items it
JOIN ingredients i ON it.ingredient_id = i.id
WHERE it.grocery_list_id = ?;

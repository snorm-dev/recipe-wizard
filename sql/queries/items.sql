-- name: CreateItem :one
INSERT INTO items (created_at, updated_at, ingredient_id, grocery_list_id, meal_id, name, description, amount, units, standard_amount, standard_units)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = ?;

-- name: GetItemAndGroceryList :one
SELECT sqlc.embed(it), sqlc.embed(gl) FROM items it
JOIN grocery_lists gl ON it.grocery_list_id = gl.id
WHERE it.id = ?;

-- name: GetItemsForMeal :many
SELECT * FROM items it 
WHERE it.meal_id = ?;

-- name: GetExtendedItem :one
SELECT sqlc.embed(it), sqlc.embed(i) FROM items it
LEFT JOIN ingredients i ON it.ingredient_id = i.id
WHERE it.id = ?;

-- name: GetExtendedItemsForMeal :many
SELECT sqlc.embed(it), sqlc.embed(i) FROM items it
LEFT JOIN ingredients i ON it.ingredient_id = i.id
WHERE it.meal_id = ?;

-- name: GetItemsForGroceryList :many
SELECT * FROM items it 
WHERE it.grocery_list_id = ?;

-- name: GetExtendedItemsForGroceryList :many
SELECT sqlc.embed(it), sqlc.embed(i) FROM items it
LEFT JOIN ingredients i ON it.ingredient_id = i.id
WHERE it.grocery_list_id = ?;

-- name: GetItemsForGroceryListByName :many
SELECT * FROM items it 
WHERE it.grocery_list_id = ? AND name = ?;

-- name: SetIsComplete :exec
UPDATE items
SET updated_at = ?, is_complete = ?
WHERE id = ?;

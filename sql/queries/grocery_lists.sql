-- name: CreateGroceryList :one
INSERT INTO grocery_lists (created_at, updated_at, name, owner_id)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetGroceryList :one
SELECT * FROM grocery_lists
WHERE id = ?;

-- name: GetGroceryListsForUser :many
SELECT * FROM grocery_lists
WHERE owner_id = ?;

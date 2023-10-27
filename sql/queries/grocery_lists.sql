-- name: CreateGroceryList :exec
INSERT INTO grocery_list (created_at, updated_at, name, owner_id)
VALUES (?, ?, ?, ?);

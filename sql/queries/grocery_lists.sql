-- name: CreateGroceryList :execresult
INSERT INTO grocery_lists (created_at, updated_at, name, owner_id)
VALUES (?, ?, ?, ?);

-- name: GetGroceryList :one
SELECT * FROM grocery_lists
WHERE id = ?;

-- name: GetGroceryListsForUser :many
SELECT * FROM grocery_lists
WHERE owner_id = ?;

-- name: GetIngredientsInGroceryList :many
SELECT i.* FROM ingredients i
JOIN recipes r ON r.id = i.recipe_id
JOIN recipe_instances ri ON r.id = ri.recipe_id
WHERE ri.grocery_list_id = ?;

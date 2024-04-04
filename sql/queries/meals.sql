-- name: CreateMeal :one
INSERT INTO meals (created_at, updated_at, grocery_list_id, recipe_id)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetMeal :one
SELECT * FROM meals
WHERE id = ?;

-- name: GetMealsInGroceryList :many
SELECT * FROM meals m 
WHERE m.grocery_list_id = ?;

-- name: GetExtendedMeal :one
SELECT sqlc.embed(m), sqlc.embed(r) from meals m 
JOIN recipes r ON m.recipe_id = r.id
WHERE m.id = ?;

-- name: GetExtendedMealsInGroceryList :many
SELECT sqlc.embed(m), sqlc.embed(r) from meals m 
JOIN recipes r ON m.recipe_id = r.id
WHERE m.grocery_list_id = ?;

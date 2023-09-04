// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: recipe_ingredients.sql

package database

import (
	"context"
	"time"
)

const createRecipeIngredient = `-- name: CreateRecipeIngredient :exec
INSERT INTO recipe_ingredients(id, created_at, updated_at, ingredient_id, recipe_id, quantity, units)
VALUES (?, ?, ?, ?, ?, ?, ?)
`

type CreateRecipeIngredientParams struct {
	ID           string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IngredientID string
	RecipeID     string
	Quantity     int32
	Units        string
}

func (q *Queries) CreateRecipeIngredient(ctx context.Context, arg CreateRecipeIngredientParams) error {
	_, err := q.db.ExecContext(ctx, createRecipeIngredient,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.IngredientID,
		arg.RecipeID,
		arg.Quantity,
		arg.Units,
	)
	return err
}

const getIngredientsForRecipe = `-- name: GetIngredientsForRecipe :many
SELECT ri.id, ri.created_at, ri.updated_at, ri.ingredient_id, ri.recipe_id, ri.quantity, ri.units, i.id, i.created_at, i.updated_at, i.name, i.description FROM recipe_ingredients ri
LEFT JOIN ingredients i ON i.id = ri.ingredient_id
WHERE ri.recipe_id = ?
`

type GetIngredientsForRecipeRow struct {
	RecipeIngredient RecipeIngredient
	Ingredient       Ingredient
}

func (q *Queries) GetIngredientsForRecipe(ctx context.Context, recipeID string) ([]GetIngredientsForRecipeRow, error) {
	rows, err := q.db.QueryContext(ctx, getIngredientsForRecipe, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetIngredientsForRecipeRow
	for rows.Next() {
		var i GetIngredientsForRecipeRow
		if err := rows.Scan(
			&i.RecipeIngredient.ID,
			&i.RecipeIngredient.CreatedAt,
			&i.RecipeIngredient.UpdatedAt,
			&i.RecipeIngredient.IngredientID,
			&i.RecipeIngredient.RecipeID,
			&i.RecipeIngredient.Quantity,
			&i.RecipeIngredient.Units,
			&i.Ingredient.ID,
			&i.Ingredient.CreatedAt,
			&i.Ingredient.UpdatedAt,
			&i.Ingredient.Name,
			&i.Ingredient.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecipeIngredient = `-- name: GetRecipeIngredient :one
SELECT id, created_at, updated_at, ingredient_id, recipe_id, quantity, units FROM recipe_ingredients
WHERE id = ?
`

func (q *Queries) GetRecipeIngredient(ctx context.Context, id string) (RecipeIngredient, error) {
	row := q.db.QueryRowContext(ctx, getRecipeIngredient, id)
	var i RecipeIngredient
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IngredientID,
		&i.RecipeID,
		&i.Quantity,
		&i.Units,
	)
	return i, err
}
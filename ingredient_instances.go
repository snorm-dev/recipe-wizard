package main

import (
	"time"

	"github.com/snorman7384/recipe-wizard/internal/database"
)

type ingredientInstanceResponse struct {
	ID               int64     `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	GroceryListID    int64     `json:"grocery_list_id"`
	RecipeInstanceID *int64    `json:"recipe_instance_id,omitempty"`
	IngredientID     *int64    `json:"ingredient_id,omitempty"`
}

func databaseIngredientInstanceToResponse(ii database.IngredientInstance) ingredientInstanceResponse {
	return ingredientInstanceResponse{
		ID:               ii.ID,
		CreatedAt:        ii.CreatedAt,
		UpdatedAt:        ii.UpdatedAt,
		GroceryListID:    ii.GroceryListID,
		RecipeInstanceID: int64PointerFromSqlNullInt64(ii.RecipeInstanceID),
		IngredientID:     int64PointerFromSqlNullInt64(ii.IngredientID),
	}
}

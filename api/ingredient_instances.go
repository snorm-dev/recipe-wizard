package api

import (
	"time"

	"github.com/snorman7384/recipe-wizard/internal/database"
)

type ingredientInstanceResponse struct {
	ID               int64     `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	GroceryListID    int64     `json:"grocery_list_id"`
	RecipeInstanceID int64     `json:"recipe_instance_id,omitempty"` // 0 is never a sql id, so we can treat 0 as "no recipe instance"
	IngredientID     int64     `json:"ingredient_id,omitempty"`      // likewise
}

func databaseIngredientInstanceToResponse(ii database.IngredientInstance) ingredientInstanceResponse {
	return ingredientInstanceResponse{
		ID:               ii.ID,
		CreatedAt:        ii.CreatedAt,
		UpdatedAt:        ii.UpdatedAt,
		GroceryListID:    ii.GroceryListID,
		RecipeInstanceID: ii.RecipeInstanceID.Int64,
		IngredientID:     ii.IngredientID.Int64,
	}
}

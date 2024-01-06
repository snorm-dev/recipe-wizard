package domain

import (
	"context"
	"database/sql"

	"github.com/snorman7384/recipe-wizard/internal/database"
)

func (c *Config) GetIngredientInstancesForRecipeInstance(ctx context.Context, recipeInstance database.RecipeInstance) ([]database.IngredientInstance, error) {

	ingredientInstances, err := c.Querier().GetIngredientInstancesForRecipeInstance(ctx, sql.NullInt64{Valid: true, Int64: recipeInstance.ID})
	if err != nil {
		return nil, err
	}

	return ingredientInstances, nil
}

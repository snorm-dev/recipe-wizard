package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/snorman7384/recipe-wizard/internal/database"
)

func (c *Config) CreateRecipeInstance(ctx context.Context, user database.User, groceryList database.GroceryList, recipeID int64) (database.RecipeInstance, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return database.RecipeInstance{}, err
	}
	defer tx.Rollback()

	qtx := c.Querier().WithTx(tx)

	now := time.Now()

	result, err := qtx.CreateRecipeInstance(ctx, database.CreateRecipeInstanceParams{
		CreatedAt:     now,
		UpdatedAt:     now,
		RecipeID:      recipeID,
		GroceryListID: groceryList.ID,
	})
	if err != nil {
		return database.RecipeInstance{}, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return database.RecipeInstance{}, nil
	}

	recipeInstance, err := qtx.GetRecipeInstance(ctx, id)
	if err != nil {
		return database.RecipeInstance{}, nil
	}

	ingredients, err := qtx.GetIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return database.RecipeInstance{}, nil
	}

	for _, ingredient := range ingredients {

		now := time.Now()

		_, err := qtx.CreateIngredientInstance(ctx, database.CreateIngredientInstanceParams{
			CreatedAt:        now,
			UpdatedAt:        now,
			IngredientID:     sql.NullInt64{Int64: ingredient.ID, Valid: true},
			GroceryListID:    groceryList.ID,
			RecipeInstanceID: sql.NullInt64{Int64: recipeInstance.ID, Valid: true},
		})
		if err != nil {
			return database.RecipeInstance{}, nil
		}
	}

	return recipeInstance, tx.Commit()
}

func (c *Config) GetRecipeInstancesInGroceryList(ctx context.Context, groceryList database.GroceryList) ([]database.RecipeInstance, error) {

	recipeInstances, err := c.Querier().GetRecipeInstancesInGroceryList(ctx, groceryList.ID)
	if err != nil {
		return nil, err
	}

	return recipeInstances, nil
}

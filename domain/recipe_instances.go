package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainRecipeInstance(ri database.RecipeInstance, recipe Recipe) RecipeInstance {
	return RecipeInstance{
		ID:            ri.ID,
		CreatedAt:     ri.CreatedAt,
		UpdatedAt:     ri.UpdatedAt,
		GroceryListID: ri.GroceryListID,
		Recipe:        recipe,
	}
}

func (c *Config) CreateRecipeInstance(ctx context.Context, user User, groceryList GroceryList, recipeID int64) (RecipeInstance, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return RecipeInstance{}, err
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
		return RecipeInstance{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return RecipeInstance{}, err
	}

	row, err := qtx.GetExtendedRecipeInstance(ctx, id)
	if err != nil {
		return RecipeInstance{}, err
	}

	ingredients, err := qtx.GetIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return RecipeInstance{}, err
	}

	for _, ingredient := range ingredients {

		now := time.Now()

		_, err := qtx.CreateItem(ctx, database.CreateItemParams{
			CreatedAt:        now,
			UpdatedAt:        now,
			IngredientID:     sql.NullInt64{Int64: ingredient.ID, Valid: true},
			GroceryListID:    groceryList.ID,
			RecipeInstanceID: sql.NullInt64{Int64: row.RecipeInstance.ID, Valid: true},
		})
		if err != nil {
			return RecipeInstance{}, err
		}
	}

	return databaseToDomainRecipeInstance(row.RecipeInstance, databaseToDomainRecipe(row.Recipe)), tx.Commit()
}

func (c *Config) GetRecipeInstancesInGroceryList(ctx context.Context, groceryList GroceryList) ([]RecipeInstance, error) {

	rows, err := c.Querier().GetExtendedRecipeInstancesInGroceryList(ctx, groceryList.ID)
	if err != nil {
		return nil, err
	}

	recipeInstances := make([]RecipeInstance, len(rows))

	for i, row := range rows {
		recipeInstances[i] = databaseToDomainRecipeInstance(row.RecipeInstance, databaseToDomainRecipe(row.Recipe))
	}

	return recipeInstances, nil
}

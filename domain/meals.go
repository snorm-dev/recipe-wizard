package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainMeal(m database.Meal, recipe Recipe) Meal {
	return Meal{
		ID:            m.ID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		GroceryListID: m.GroceryListID,
		Recipe:        recipe,
	}
}

func (c *Config) CreateMeal(ctx context.Context, user User, groceryList GroceryList, recipeID int64) (Meal, error) {
	tx, err := c.DB.Begin()
	if err != nil {
		return Meal{}, err
	}
	defer tx.Rollback()

	qtx := c.Querier().WithTx(tx)

	now := time.Now()

	meal, err := qtx.CreateMeal(ctx, database.CreateMealParams{
		CreatedAt:     now,
		UpdatedAt:     now,
		RecipeID:      recipeID,
		GroceryListID: groceryList.ID,
	})
	if err != nil {
		return Meal{}, err
	}

	recipe, err := qtx.GetRecipe(ctx, recipeID)
	if err != nil {
		return Meal{}, err
	}

	ingredients, err := qtx.GetIngredientsForRecipe(ctx, recipeID)
	if err != nil {
		return Meal{}, err
	}

	for _, ingredient := range ingredients {

		now := time.Now()

		_, err := qtx.CreateItem(ctx, database.CreateItemParams{
			CreatedAt:      now,
			UpdatedAt:      now,
			IngredientID:   sql.NullInt64{Int64: ingredient.ID, Valid: true},
			GroceryListID:  groceryList.ID,
			MealID:         sql.NullInt64{Int64: meal.ID, Valid: true},
			Name:           ingredient.Name,
			Description:    ingredient.Description,
			Amount:         ingredient.Amount,
			Units:          ingredient.Units,
			StandardAmount: ingredient.StandardAmount,
			StandardUnits:  ingredient.StandardUnits,
		})
		if err != nil {
			return Meal{}, err
		}
	}

	return databaseToDomainMeal(meal, databaseToDomainRecipe(recipe)), tx.Commit()
}

func (c *Config) GetMealsInGroceryList(ctx context.Context, groceryList GroceryList) ([]Meal, error) {

	rows, err := c.Querier().GetExtendedMealsInGroceryList(ctx, groceryList.ID)
	if err != nil {
		return nil, err
	}

	meals := make([]Meal, len(rows))

	for i, row := range rows {
		meals[i] = databaseToDomainMeal(row.Meal, databaseToDomainRecipe(row.Recipe))
	}

	return meals, nil
}

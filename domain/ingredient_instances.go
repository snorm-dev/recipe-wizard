package domain

import (
	"context"
	"database/sql"

	"github.com/snorman7384/recipe-wizard/ingparse"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainIngredientInstance(ii database.IngredientInstance, ingredient Ingredient) IngredientInstance {
	return IngredientInstance{
		ID:               ii.ID,
		CreatedAt:        ii.CreatedAt,
		UpdatedAt:        ii.UpdatedAt,
		GroceryListID:    ii.GroceryListID,
		RecipeInstanceID: ii.RecipeInstanceID.Int64,
		Ingredient:       ingredient,
	}
}

func (c *Config) GetIngredientInstancesForRecipeInstance(ctx context.Context, recipeInstance RecipeInstance) ([]IngredientInstance, error) {

	rows, err := c.Querier().GetExtendedIngredientInstancesForRecipeInstance(ctx, sql.NullInt64{Valid: true, Int64: recipeInstance.ID})
	if err != nil {
		return nil, err
	}

	ingredientInstances := make([]IngredientInstance, len(rows))

	for i, row := range rows {
		ingredientInstances[i] = databaseToDomainIngredientInstance(row.IngredientInstance, databaseToDomainIngredient(row.Ingredient))
	}

	return ingredientInstances, nil
}

func (c *Config) GetIngredientInstancesForGroceryList(ctx context.Context, groceryList GroceryList) ([]IngredientInstance, error) {
	rows, err := c.Querier().GetExtendedIngredientInstancesForGroceryList(ctx, groceryList.ID)
	if err != nil {
		return nil, err
	}

	ingredientInstances := make([]IngredientInstance, len(rows))

	for i, row := range rows {
		ingredientInstances[i] = databaseToDomainIngredientInstance(row.IngredientInstance, databaseToDomainIngredient(row.Ingredient))
	}

	return ingredientInstances, nil
}

// TODO
/* func (c *Config) GetIngredientInstancesForGroceryListByName(ctx context.Context, groceryList GroceryList, name string) ([]IngredientInstance, error) {
	// more complicated sql call
	return nil, nil
} */

func (c *Config) GetIngredientGroupsForGroceryList(ctx context.Context, groceryList GroceryList) ([]IngredientGroup, error) {
	ingredientInstances, err := c.GetIngredientInstancesForGroceryList(ctx, groceryList)
	if err != nil {
		return nil, err
	}

	totalsMap := make(map[string]map[ingparse.StandardUnit]IngredientGroup)

	for _, ii := range ingredientInstances {

		name := ii.Ingredient.Name
		if _, ok := totalsMap[name]; !ok {
			totalsMap[name] = make(map[ingparse.StandardUnit]IngredientGroup)
		}

		units := ii.Ingredient.StandardUnits
		if entry, ok := totalsMap[name][units]; !ok {
			totalsMap[name][units] = IngredientGroup{
				Name:      name,
				Units:     units,
				Total:     ii.Ingredient.StandardAmount,
				Instances: append(make([]IngredientInstance, 0), ii),
			}
		} else {
			entry.Total += ii.Ingredient.StandardAmount
			entry.Instances = append(entry.Instances, ii)
			totalsMap[name][units] = entry
		}

	}

	groups := make([]IngredientGroup, 0)

	for _, unitsMap := range totalsMap {
		for _, group := range unitsMap {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

// TODO
/* func (c *Config) GetIngredientGroupForGroceryListByName(ctx context.Context, groceryList GroceryList, name string) ([]IngredientGroup, error) {
	ingredients, err := c.GetIngredientInstancesForGroceryListByName(ctx, groceryList, name)
	if err != nil {
		return nil, err
	}

	unitsMap := make(map[ingparse.StandardUnit]IngredientGroup)

	for _, ingredient := range ingredients {

	}

	return nil, nil
} */

package domain

import (
	"context"
	"database/sql"

	"github.com/snorman7384/recipe-wizard/ingparse"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainItem(it database.Item, ingredient Ingredient) Item {
	return Item{
		ID:               it.ID,
		CreatedAt:        it.CreatedAt,
		UpdatedAt:        it.UpdatedAt,
		GroceryListID:    it.GroceryListID,
		RecipeInstanceID: it.RecipeInstanceID.Int64,
		Ingredient:       ingredient,
		Name:             it.Name,
		Description:      it.Description.String,
		Amount:           it.Amount,
		Units:            it.Units,
		StandardAmount:   it.StandardAmount,
		StandardUnits:    ingparse.StandardUnitFromString(it.StandardUnits),
	}
}

func (c *Config) GetItemsForRecipeInstance(ctx context.Context, recipeInstance RecipeInstance) ([]Item, error) {

	rows, err := c.Querier().GetExtendedItemsForRecipeInstance(ctx, sql.NullInt64{Valid: true, Int64: recipeInstance.ID})
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(rows))

	for i, row := range rows {
		items[i] = databaseToDomainItem(row.Item, databaseToDomainIngredient(row.Ingredient))
	}

	return items, nil
}

func (c *Config) GetItemsForGroceryList(ctx context.Context, groceryList GroceryList) ([]Item, error) {
	rows, err := c.Querier().GetExtendedItemsForGroceryList(ctx, groceryList.ID)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(rows))

	for i, row := range rows {
		items[i] = databaseToDomainItem(row.Item, databaseToDomainIngredient(row.Ingredient))
	}

	return items, nil
}

// TODO
/* func (c *Config) GetItemsForGroceryListByName(ctx context.Context, groceryList GroceryList, name string) ([]Item, error) {
	// more complicated sql call
	return nil, nil
} */

func (c *Config) GetItemGroupsForGroceryList(ctx context.Context, groceryList GroceryList) ([]ItemGroup, error) {
	items, err := c.GetItemsForGroceryList(ctx, groceryList)
	if err != nil {
		return nil, err
	}

	totalsMap := make(map[string]map[ingparse.StandardUnit]ItemGroup)

	for _, it := range items {

		name := it.Ingredient.Name
		if _, ok := totalsMap[name]; !ok {
			totalsMap[name] = make(map[ingparse.StandardUnit]ItemGroup)
		}

		units := it.Ingredient.StandardUnits
		if entry, ok := totalsMap[name][units]; !ok {
			totalsMap[name][units] = ItemGroup{
				Name:  name,
				Units: units,
				Total: it.Ingredient.StandardAmount,
				Items: append(make([]Item, 0), it),
			}
		} else {
			entry.Total += it.Ingredient.StandardAmount
			entry.Items = append(entry.Items, it)
			totalsMap[name][units] = entry
		}

	}

	groups := make([]ItemGroup, 0)

	for _, unitsMap := range totalsMap {
		for _, group := range unitsMap {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

// TODO
/* func (c *Config) GetItemGroupForGroceryListByName(ctx context.Context, groceryList GroceryList, name string) ([]IngredientGroup, error) {
	ingredients, err := c.GetItemsForGroceryListByName(ctx, groceryList, name)
	if err != nil {
		return nil, err
	}

	unitsMap := make(map[ingparse.StandardUnit]IngredientGroup)

	for _, ingredient := range ingredients {

	}

	return nil, nil
} */

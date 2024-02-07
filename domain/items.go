package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/ingparse"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainItem(it database.Item) Item {
	return Item{
		ID:               it.ID,
		CreatedAt:        it.CreatedAt,
		UpdatedAt:        it.UpdatedAt,
		GroceryListID:    it.GroceryListID,
		RecipeInstanceID: it.RecipeInstanceID.Int64,
		IngredientID:     it.IngredientID.Int64,
		Name:             it.Name,
		Description:      it.Description.String,
		Amount:           it.Amount,
		Units:            it.Units,
		StandardAmount:   it.StandardAmount,
		StandardUnits:    ingparse.StandardUnitFromString(it.StandardUnits),
	}
}

func (c *Config) CreateItem(ctx context.Context, groceryList GroceryList, name string, description string, amount float64, units string) (Item, error) {
	now := time.Now()

	tx, err := c.DB.Begin()
	if err != nil {
		return Item{}, err
	}
	defer tx.Rollback()

	qtx := c.Querier().WithTx(tx)

	result, err := qtx.CreateItem(ctx, database.CreateItemParams{
		CreatedAt:        now,
		UpdatedAt:        now,
		IngredientID:     sql.NullInt64{},
		RecipeInstanceID: sql.NullInt64{},
		GroceryListID:    groceryList.ID,
		Name:             name,
		Description:      sql.NullString{String: description, Valid: description == ""},
		Amount:           amount,
		Units:            units,
		StandardAmount:   -1,      // TODO
		StandardUnits:    "whole", // TODO
	})
	if err != nil {
		return Item{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Item{}, err
	}

	item, err := qtx.GetItem(ctx, id)
	if err != nil {
		return Item{}, err
	}

	return databaseToDomainItem(item), tx.Commit()
}

func (c *Config) GetItem(ctx context.Context, user User, id int64) (Item, error) {
	row, err := c.Querier().GetItemAndGroceryList(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Item{}, domerr.ErrNotFound
	}
	if err != nil {
		return Item{}, err
	}

	if user.ID != row.GroceryList.OwnerID {
		return Item{}, domerr.ErrForbidden
	}

	return databaseToDomainItem(row.Item), nil
}

func (c *Config) GetItemsForRecipeInstance(ctx context.Context, recipeInstance RecipeInstance) ([]Item, error) {

	dbItems, err := c.Querier().GetItemsForRecipeInstance(ctx, sql.NullInt64{Valid: true, Int64: recipeInstance.ID})
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(dbItems))

	for i, dbItem := range dbItems {
		items[i] = databaseToDomainItem(dbItem)
	}

	return items, nil
}

func (c *Config) GetItemsForGroceryList(ctx context.Context, groceryList GroceryList) ([]Item, error) {
	dbItems, err := c.Querier().GetItemsForGroceryList(ctx, groceryList.ID)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(dbItems))

	for i, dbItem := range dbItems {
		items[i] = databaseToDomainItem(dbItem)
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

		name := it.Name
		if _, ok := totalsMap[name]; !ok {
			totalsMap[name] = make(map[ingparse.StandardUnit]ItemGroup)
		}

		units := it.StandardUnits
		if entry, ok := totalsMap[name][units]; !ok {
			totalsMap[name][units] = ItemGroup{
				Name:  name,
				Units: units,
				Total: it.StandardAmount,
				Items: append(make([]Item, 0), it),
			}
		} else {
			entry.Total += it.StandardAmount
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

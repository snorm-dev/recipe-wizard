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
	var status ItemStatus
	if it.IsComplete {
		status = Complete
	} else {
		status = Incomplete
	}
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
		Status:           status,
	}
}

func (c *Config) CreateItem(ctx context.Context, groceryList GroceryList, name string, description string, amount float64, units string) (Item, error) {
	now := time.Now()

	item, err := c.Querier().CreateItem(ctx, database.CreateItemParams{
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

	return databaseToDomainItem(item), nil
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

func (c *Config) GetItemsForGroceryListByName(ctx context.Context, groceryList GroceryList, name string) ([]Item, error) {
	dbItems, err := c.Querier().GetItemsForGroceryListByName(ctx, database.GetItemsForGroceryListByNameParams{
		GroceryListID: groceryList.ID,
		Name:          name,
	})
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(dbItems))

	for i, dbItem := range dbItems {
		items[i] = databaseToDomainItem(dbItem)
	}

	return items, nil
}

func (c *Config) GetItemGroupsForGroceryList(ctx context.Context, groceryList GroceryList) ([]ItemGroup, error) {
	items, err := c.GetItemsForGroceryList(ctx, groceryList)
	if err != nil {
		return nil, err
	}

	groupMap := make(map[string]ItemGroup)

	for _, it := range items {
		name := it.Name
		entry, ok := groupMap[name]
		if !ok {
			entry = ItemGroup{
				Name:   name,
				Totals: make(map[ingparse.StandardUnit]float64),
				Items:  make([]Item, 0),
			}
		}

		entry.Totals[it.StandardUnits] += it.StandardAmount
		entry.Items = append(entry.Items, it)
		groupMap[name] = entry
	}

	groups := make([]ItemGroup, 0)

	for _, group := range groupMap {
		groups = append(groups, group)
	}

	return groups, nil
}

func (c *Config) GetItemGroupForGroceryListByName(ctx context.Context, groceryList GroceryList, name string) (ItemGroup, error) {
	items, err := c.GetItemsForGroceryListByName(ctx, groceryList, name)
	if err != nil {
		return ItemGroup{}, err
	}
	totals := make(map[ingparse.StandardUnit]float64)

	for _, item := range items {
		totals[item.StandardUnits] += item.StandardAmount
	}

	group := ItemGroup{
		Name:   name,
		Items:  items,
		Totals: totals,
	}

	return group, nil
}

func (c *Config) MarkItemStatus(ctx context.Context, item Item, status ItemStatus) (Item, error) {
	if item.Status == status {
		return item, nil
	}

	now := time.Now()
	isComplete := status == Complete

	err := c.Querier().SetIsComplete(ctx, database.SetIsCompleteParams{
		UpdatedAt:  now,
		IsComplete: isComplete,
		ID:         item.ID,
	})
	if err != nil {
		return Item{}, err
	}

	item.Status = status

	return item, nil
}

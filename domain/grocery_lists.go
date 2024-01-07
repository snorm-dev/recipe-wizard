package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainGroceryList(dbGroceryList database.GroceryList) GroceryList {
	return GroceryList{
		ID:        dbGroceryList.ID,
		CreatedAt: dbGroceryList.CreatedAt,
		UpdatedAt: dbGroceryList.UpdatedAt,
		Name:      dbGroceryList.Name,
		OwnerID:   dbGroceryList.ID,
	}
}

func (c *Config) CreateGroceryList(ctx context.Context, user User, name string) (GroceryList, error) {

	tx, err := c.DB.Begin()
	if err != nil {
		return GroceryList{}, err
	}
	defer tx.Rollback()

	qtx := c.Querier().WithTx(tx)

	now := time.Now()

	result, err := qtx.CreateGroceryList(ctx, database.CreateGroceryListParams{
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		OwnerID:   user.ID,
	})
	if err != nil {
		return GroceryList{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return GroceryList{}, err
	}

	groceryList, err := qtx.GetGroceryList(ctx, id)
	if err != nil {
		return GroceryList{}, err
	}

	return databaseToDomainGroceryList(groceryList), tx.Commit()
}

func (c *Config) GetGroceryList(ctx context.Context, user User, id int64) (GroceryList, error) {

	groceryList, err := c.Querier().GetGroceryList(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return GroceryList{}, domerr.ErrNotFound
	}
	if err != nil {
		return GroceryList{}, err
	}

	if user.ID != groceryList.OwnerID {
		return GroceryList{}, domerr.ErrForbidden
	}

	return databaseToDomainGroceryList(groceryList), nil
}

func (c *Config) GetGroceryListsForUser(ctx context.Context, user User) ([]GroceryList, error) {
	groceryLists, err := c.Querier().GetGroceryListsForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	domainList := make([]GroceryList, len(groceryLists))
	for i, groceryList := range groceryLists {
		domainList[i] = databaseToDomainGroceryList(groceryList)
	}

	return domainList, nil
}

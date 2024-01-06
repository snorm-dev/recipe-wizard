package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func (c *Config) CreateGroceryList(ctx context.Context, user database.User, name string) (database.GroceryList, error) {

	tx, err := c.DB.Begin()
	if err != nil {
		return database.GroceryList{}, err
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
		return database.GroceryList{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return database.GroceryList{}, err
	}

	groceryList, err := qtx.GetGroceryList(ctx, id)
	if err != nil {
		return database.GroceryList{}, err
	}

	return groceryList, tx.Commit()
}

func (c *Config) GetGroceryList(ctx context.Context, user database.User, id int64) (database.GroceryList, error) {

	groceryList, err := c.Querier().GetGroceryList(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return database.GroceryList{}, domerr.ErrNotFound
	}
	if err != nil {
		return database.GroceryList{}, err
	}

	if user.ID != groceryList.OwnerID {
		return database.GroceryList{}, domerr.ErrForbidden
	}

	return groceryList, nil
}

func (c *Config) GetGroceryListsForUser(ctx context.Context, user database.User) ([]database.GroceryList, error) {
	groceryLists, err := c.Querier().GetGroceryListsForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return groceryLists, nil
}

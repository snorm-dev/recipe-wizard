package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type CreateUserParams struct {
	Username       string
	HashedPassword string
	FirstName      string
	LastName       string
}

func (c *Config) CreateUser(ctx context.Context, params CreateUserParams) (database.User, error) {

	tx, err := c.DB.Begin()
	if err != nil {
		return database.User{}, err
	}
	defer tx.Rollback()

	qtx := c.Querier().WithTx(tx)

	now := time.Now()
	firstName := sql.NullString{String: params.FirstName, Valid: params.FirstName != ""}
	lastName := sql.NullString{String: params.LastName, Valid: params.LastName != ""}

	result, err := qtx.CreateUser(ctx, database.CreateUserParams{
		CreatedAt:      now,
		UpdatedAt:      now,
		Username:       params.Username,
		HashedPassword: params.HashedPassword,
		FirstName:      firstName,
		LastName:       lastName,
	})

	if err != nil {
		return database.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return database.User{}, err
	}

	user, err := qtx.GetUser(ctx, id)
	if err != nil {
		return database.User{}, err
	}

	return user, tx.Commit()
}

func (c *Config) GetUser(ctx context.Context, id int64) (database.User, error) {

	user, err := c.Querier().GetUser(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return database.User{}, domerr.ErrUserNotFound
	} else if err != nil {
		return database.User{}, err
	}

	return user, nil
}

func (c *Config) GetUserByUsername(ctx context.Context, username string) (database.User, error) {

	user, err := c.Querier().GetUserByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return database.User{}, domerr.ErrUserNotFound
	} else if err != nil {
		return database.User{}, err
	}

	return user, nil
}

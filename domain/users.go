package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainUser(user database.User) User {
	return User{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FirstName:      user.FirstName.String,
		LastName:       user.LastName.String,
	}
}

type CreateUserParams struct {
	Username       string
	HashedPassword string
	FirstName      string
	LastName       string
}

func (c *Config) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {

	now := time.Now()
	firstName := sql.NullString{String: params.FirstName, Valid: params.FirstName != ""}
	lastName := sql.NullString{String: params.LastName, Valid: params.LastName != ""}

	user, err := c.Querier().CreateUser(ctx, database.CreateUserParams{
		CreatedAt:      now,
		UpdatedAt:      now,
		Username:       params.Username,
		HashedPassword: params.HashedPassword,
		FirstName:      firstName,
		LastName:       lastName,
	})

	if err != nil {
		return User{}, err
	}

	return databaseToDomainUser(user), nil
}

func (c *Config) GetUser(ctx context.Context, id int64) (User, error) {

	user, err := c.Querier().GetUser(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, domerr.ErrUserNotFound
	} else if err != nil {
		return User{}, err
	}

	return databaseToDomainUser(user), nil
}

func (c *Config) GetUserByUsername(ctx context.Context, username string) (User, error) {

	user, err := c.Querier().GetUserByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, domerr.ErrUserNotFound
	} else if err != nil {
		return User{}, err
	}

	return databaseToDomainUser(user), nil
}

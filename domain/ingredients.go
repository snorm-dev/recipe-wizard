package domain

import (
	"context"

	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func (c *Config) GetIngredientsForRecipe(ctx context.Context, user database.User, recipe database.Recipe) ([]database.Ingredient, error) {
	if user.ID != recipe.OwnerID {
		return nil, domerr.ErrForbidden
	}

	ingredients, err := c.Querier().GetIngredientsForRecipe(ctx, recipe.ID)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

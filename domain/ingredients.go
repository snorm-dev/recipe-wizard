package domain

import (
	"context"

	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/ingparse"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func databaseToDomainIngredient(ingredient database.Ingredient) Ingredient {
	return Ingredient{
		ID:             ingredient.ID,
		CreatedAt:      ingredient.CreatedAt,
		UpdatedAt:      ingredient.UpdatedAt,
		Name:           ingredient.Name,
		Description:    ingredient.Description.String,
		Units:          ingredient.Units,
		Amount:         ingredient.Amount,
		StandardUnits:  ingparse.StandardUnitFromString(ingredient.StandardUnits),
		StandardAmount: ingredient.StandardAmount,
		RecipeID:       ingredient.RecipeID,
	}
}

func (c *Config) GetIngredientsForRecipe(ctx context.Context, user User, recipe Recipe) ([]Ingredient, error) {
	if user.ID != recipe.OwnerID {
		return nil, domerr.ErrForbidden
	}

	ingredients, err := c.Querier().GetIngredientsForRecipe(ctx, recipe.ID)
	if err != nil {
		return nil, err
	}

	domainList := make([]Ingredient, len(ingredients))

	for i, ingredient := range ingredients {
		domainList[i] = databaseToDomainIngredient(ingredient)
	}

	return domainList, nil
}

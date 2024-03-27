package domain

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/kkyr/go-recipe/pkg/recipe"
	"github.com/snorman7384/recipe-wizard/domerr"
	"github.com/snorman7384/recipe-wizard/internal/database"
	"github.com/snorman7384/recipe-wizard/internal/misc"
)

func databaseToDomainRecipe(recipe database.Recipe) Recipe {
	return Recipe{
		ID:          recipe.ID,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
		Name:        recipe.Name,
		Description: recipe.Description.String,
		Url:         recipe.Url.String,
		PrepTime:    recipe.PrepTime.String,
		CookTime:    recipe.CookTime.String,
		TotalTime:   recipe.TotalTime.String,
		OwnerID:     recipe.OwnerID,
	}
}

func (c *Config) CreateRecipeFromUrl(ctx context.Context, user User, url string) (Recipe, error) {

	// get recipe data
	s, err := recipe.ScrapeURL(url)
	if err != nil {
		// respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse recipe from url: %v", err))
		return Recipe{}, domerr.ErrRecipeScraperFailure
	}

	name, ok := s.Name()
	if !ok {
		name = url
	}

	timeDurationToString := func(t time.Duration, ok bool) (string, bool) {
		return t.String(), ok
	}

	description := misc.SqlNullStringFromOkString(s.Description())
	prepTime := misc.SqlNullStringFromOkString(timeDurationToString(s.PrepTime()))
	cookTime := misc.SqlNullStringFromOkString(timeDurationToString(s.CookTime()))
	totalTime := misc.SqlNullStringFromOkString(timeDurationToString(s.TotalTime()))

	now := time.Now()
	sqlUrl := misc.SqlNullStringFromOkString(url, true)

	tx, err := c.DB.Begin()
	if err != nil {
		return Recipe{}, err
	}
	defer tx.Rollback()

	qtx := c.Querier().WithTx(tx)

	recipe, err := qtx.CreateRecipe(ctx, database.CreateRecipeParams{
		CreatedAt:   now,
		UpdatedAt:   now,
		Url:         sqlUrl,
		Name:        name,
		Description: description,
		CookTime:    cookTime,
		PrepTime:    prepTime,
		TotalTime:   totalTime,
		OwnerID:     user.ID,
	})
	if err != nil {
		return Recipe{}, err
	}

	if ings, ok := s.Ingredients(); ok {
		ingredients, err := c.IngredientParser.ParseIngredients(ings)
		if err != nil {
			return Recipe{}, err
		}
		var wg sync.WaitGroup
		ch := make(chan error, len(ingredients))
		for _, ingredient := range ingredients {
			ingredient := ingredient // I love loop variables

			wg.Add(1)

			go func() {
				defer wg.Done()

				now = time.Now()
				_, err := qtx.CreateIngredient(ctx, database.CreateIngredientParams{
					CreatedAt:      now,
					UpdatedAt:      now,
					Name:           ingredient.Name,
					Amount:         ingredient.Measure.OriginalAmount,
					Units:          ingredient.Measure.OriginalUnits,
					StandardAmount: ingredient.Measure.StandardAmount,
					StandardUnits:  ingredient.Measure.StandardUnits.String(),
					RecipeID:       recipe.ID,
					Description:    sql.NullString{String: ingredient.Description, Valid: ingredient.Description != ""},
				})
				if err != nil {
					ch <- err
					return
				}
			}()
		}
		wg.Wait()
		select {
		case err := <-ch:
			return Recipe{}, err
		default:
		}
	}

	return databaseToDomainRecipe(recipe), tx.Commit()
}

func (c *Config) GetRecipe(ctx context.Context, user User, id int64) (Recipe, error) {

	recipe, err := c.Querier().GetRecipe(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return Recipe{}, domerr.ErrNotFound
	}
	if err != nil {
		return Recipe{}, err
	}

	if user.ID != recipe.OwnerID {
		return Recipe{}, domerr.ErrForbidden
	}

	return databaseToDomainRecipe(recipe), nil
}

func (c *Config) GetRecipesForUser(ctx context.Context, user User) ([]Recipe, error) {
	recipes, err := c.Querier().GetRecipesForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	domainList := make([]Recipe, len(recipes))

	for i, recipe := range recipes {
		domainList[i] = databaseToDomainRecipe(recipe)
	}

	return domainList, nil
}

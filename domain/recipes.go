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

func (c *Config) CreateRecipeFromUrl(ctx context.Context, user database.User, url string) (database.Recipe, error) {

	// get recipe data
	s, err := recipe.ScrapeURL(url)
	if err != nil {
		// respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse recipe from url: %v", err))
		return database.Recipe{}, domerr.ErrRecipeScraperFailure
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
		return database.Recipe{}, err
	}
	defer tx.Rollback()

	qtx := c.Querier().WithTx(tx)

	result, err := qtx.CreateRecipe(ctx, database.CreateRecipeParams{
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
		return database.Recipe{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return database.Recipe{}, err
	}

	recipe, err := qtx.GetRecipe(ctx, id)
	if err != nil {
		return database.Recipe{}, err
	}

	if ings, ok := s.Ingredients(); ok {
		ingredients, err := c.IngredientParser.ParseIngredients(ings)
		if err != nil {
			return database.Recipe{}, err
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
					RecipeID:       id,
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
			return database.Recipe{}, err
		default:
		}
	}

	return recipe, tx.Commit()
}

func (c *Config) GetRecipe(ctx context.Context, user database.User, id int64) (database.Recipe, error) {

	recipe, err := c.Querier().GetRecipe(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return database.Recipe{}, domerr.ErrNotFound
	}
	if err != nil {
		return database.Recipe{}, err
	}

	if user.ID != recipe.OwnerID {
		return database.Recipe{}, domerr.ErrForbidden
	}

	return recipe, nil
}

func (c *Config) GetRecipesForUser(ctx context.Context, user database.User) ([]database.Recipe, error) {
	recipes, err := c.Querier().GetRecipesForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

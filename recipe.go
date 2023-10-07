package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kkyr/go-recipe/pkg/recipe"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func (c *config) handlePostRecipe() http.HandlerFunc {
	type request struct {
		Url string `json:"url"`
	}

	type ingredient struct {
		ID          int64     `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Name        string    `json:"name"`
		Description *string   `json:"description"`
	}

	type response struct {
		ID          int64        `json:"id"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
		Name        string       `json:"name"`
		Description *string      `json:"description,omitempty"`
		Url         *string      `json:"url,omitempty"`
		PrepTime    *string      `json:"prep_time,omitempty"`
		CookTime    *string      `json:"cook_time,omitempty"`
		TotalTime   *string      `json:"total_time,omitempty"`
		Ingredients []ingredient `json:"ingredients"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := request{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// get recipe data
		s, err := recipe.ScrapeURL(reqBody.Url)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse recipe from url: %v", err))
			return
		}

		name, ok := s.Name()
		if !ok {
			name = reqBody.Url
		}

		timeDurationToString := func(t time.Duration, ok bool) (string, bool) {
			return t.String(), ok
		}

		description := sqlNullStringFromOkString(s.Description())
		prepTime := sqlNullStringFromOkString(timeDurationToString(s.PrepTime()))
		cookTime := sqlNullStringFromOkString(timeDurationToString(s.CookTime()))
		totalTime := sqlNullStringFromOkString(timeDurationToString(s.TotalTime()))

		now := time.Now()
		url := sql.NullString{
			String: reqBody.Url,
			Valid:  true,
		}

		err = c.DB.CreateRecipe(r.Context(), database.CreateRecipeParams{
			CreatedAt:   now,
			UpdatedAt:   now,
			Url:         url,
			Name:        name,
			Description: description,
			CookTime:    cookTime,
			PrepTime:    prepTime,
			TotalTime:   totalTime,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create recipe in db: %v", err))
			return
		}

		id, err := c.DB.GetLastInsertID(r.Context())

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not retrieve id: %v", err))
		}

		recipe, err := c.DB.GetRecipe(r.Context(), id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get recipe from db: %v", err))
			return
		}

		if ingredients, ok := s.Ingredients(); ok {
			for _, ingredient := range ingredients {
				now = time.Now()
				err = c.DB.CreateIngredient(r.Context(), database.CreateIngredientParams{
					CreatedAt:   now,
					UpdatedAt:   now,
					Name:        ingredient,
					RecipeID:    id,
					Description: sql.NullString{Valid: false},
				})

				if err != nil {
					continue
				}
			}
		}

		dbIngredients, err := c.DB.GetIngredientsForRecipe(r.Context(), id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := response{
			ID:          id,
			CreatedAt:   recipe.CreatedAt,
			UpdatedAt:   recipe.UpdatedAt,
			Name:        name,
			Description: stringPointerFromSqlNullString(recipe.Description),
			Url:         stringPointerFromSqlNullString(recipe.Url),
			PrepTime:    stringPointerFromSqlNullString(recipe.PrepTime),
			CookTime:    stringPointerFromSqlNullString(recipe.CookTime),
			TotalTime:   stringPointerFromSqlNullString(recipe.TotalTime),
			Ingredients: make([]ingredient, len(dbIngredients)),
		}

		for i, dbIngredient := range dbIngredients {
			ingredient := ingredient{
				ID:        dbIngredient.ID,
				CreatedAt: dbIngredient.CreatedAt,
				UpdatedAt: dbIngredient.UpdatedAt,
				Name:      dbIngredient.Name,
			}
			ingredient.Description = stringPointerFromSqlNullString(dbIngredient.Description)

			resBody.Ingredients[i] = ingredient
		}

		respondWithJSON(w, http.StatusCreated, &resBody)
	}
}

func (c *config) handleGetRecipe() http.HandlerFunc {
	type response struct {
		ID          int64     `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Name        string    `json:"name"`
		Description *string   `json:"description"`
		Url         *string   `json:"url"`
		PrepTime    *string   `json:"prep_time"`
		CookTime    *string   `json:"cook_time"`
		TotalTime   *string   `json:"total_time"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "recipe_id")

		id, err := strconv.Atoi(idString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		recipe, err := c.DB.GetRecipe(r.Context(), int64(id))

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := response{
			ID:          recipe.ID,
			CreatedAt:   recipe.CreatedAt,
			UpdatedAt:   recipe.UpdatedAt,
			Name:        recipe.Name,
			Description: stringPointerFromSqlNullString(recipe.Description),
			Url:         stringPointerFromSqlNullString(recipe.Url),
			PrepTime:    stringPointerFromSqlNullString(recipe.PrepTime),
			CookTime:    stringPointerFromSqlNullString(recipe.CookTime),
			TotalTime:   stringPointerFromSqlNullString(recipe.TotalTime),
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

func (c *config) handleGetRecipes() http.HandlerFunc {
	type recipe struct {
		ID          int64     `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Name        string    `json:"name"`
		Description *string   `json:"description"`
		Url         *string   `json:"url"`
		PrepTime    *string   `json:"prep_time"`
		CookTime    *string   `json:"cook_time"`
		TotalTime   *string   `json:"total_time"`
	}
	type response []recipe

	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := c.DB.GetRecipes(r.Context())

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := make(response, len(recipes))

		for i, dbRecipe := range recipes {
			r := recipe{
				ID:          dbRecipe.ID,
				CreatedAt:   dbRecipe.CreatedAt,
				UpdatedAt:   dbRecipe.UpdatedAt,
				Name:        dbRecipe.Name,
				Description: stringPointerFromSqlNullString(dbRecipe.Description),
				Url:         stringPointerFromSqlNullString(dbRecipe.Url),
				PrepTime:    stringPointerFromSqlNullString(dbRecipe.PrepTime),
				CookTime:    stringPointerFromSqlNullString(dbRecipe.CookTime),
				TotalTime:   stringPointerFromSqlNullString(dbRecipe.TotalTime),
			}
			resBody[i] = r
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

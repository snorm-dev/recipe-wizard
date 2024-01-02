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

type recipeResponse struct {
	ID          int64                `json:"id"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Name        string               `json:"name"`
	Description *string              `json:"description,omitempty"`
	Url         *string              `json:"url,omitempty"`
	PrepTime    *string              `json:"prep_time,omitempty"`
	CookTime    *string              `json:"cook_time,omitempty"`
	TotalTime   *string              `json:"total_time,omitempty"`
	OwnerId     int64                `json:"owner_id"`
	Ingredients []ingredientResponse `json:"ingredients,omitempty"`
}

func databaseRecipeToResponse(recipe database.Recipe, ingredients []database.Ingredient) recipeResponse {
	responseIngredients := make([]ingredientResponse, 0, len(ingredients))
	for _, ingredient := range ingredients {
		responseIngredients = append(responseIngredients, databaseIngredientToReponse(ingredient))
	}

	return recipeResponse{
		ID:          recipe.ID,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
		Name:        recipe.Name,
		Description: stringPointerFromSqlNullString(recipe.Description),
		Url:         stringPointerFromSqlNullString(recipe.Url),
		PrepTime:    stringPointerFromSqlNullString(recipe.PrepTime),
		CookTime:    stringPointerFromSqlNullString(recipe.CookTime),
		TotalTime:   stringPointerFromSqlNullString(recipe.TotalTime),
		OwnerId:     recipe.OwnerID,
		Ingredients: responseIngredients,
	}
}

func (c *config) handlePostRecipe() http.HandlerFunc {
	type request struct {
		Url string `json:"url"`
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

		var user database.User
		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		result, err := c.DB.CreateRecipe(r.Context(), database.CreateRecipeParams{
			CreatedAt:   now,
			UpdatedAt:   now,
			Url:         url,
			Name:        name,
			Description: description,
			CookTime:    cookTime,
			PrepTime:    prepTime,
			TotalTime:   totalTime,
			OwnerID:     user.ID,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create recipe in db: %v", err))
			return
		}

		id, err := result.LastInsertId()

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
				_, err := c.DB.CreateIngredient(r.Context(), database.CreateIngredientParams{
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

		ingredients, err := c.DB.GetIngredientsForRecipe(r.Context(), id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := databaseRecipeToResponse(recipe, ingredients)

		respondWithJSON(w, http.StatusCreated, &resBody)
	}
}

func (c *config) handleGetRecipe() http.HandlerFunc {
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

		var user database.User
		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		if user.ID != recipe.OwnerID {
			respondWithError(w, http.StatusForbidden, "User does not own recipe.")
			return
		}

		ingredients, err := c.DB.GetIngredientsForRecipe(r.Context(), recipe.ID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := databaseRecipeToResponse(recipe, ingredients)

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

func (c *config) handleGetRecipes() http.HandlerFunc {
	type response []recipeResponse

	return func(w http.ResponseWriter, r *http.Request) {

		var user database.User
		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		recipes, err := c.DB.GetRecipesForUser(r.Context(), user.ID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := make(response, len(recipes))

		for i, recipe := range recipes {

			ingredients, err := c.DB.GetIngredientsForRecipe(r.Context(), recipe.ID)

			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			r := databaseRecipeToResponse(recipe, ingredients)
			resBody[i] = r
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

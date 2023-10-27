package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type ingredientResponse struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	RecipeID    int64     `json:"recipe_id"`
}

func databaseIngredientToReponse(ingredient database.Ingredient) ingredientResponse {
	return ingredientResponse{
		ID:          ingredient.ID,
		CreatedAt:   ingredient.CreatedAt,
		UpdatedAt:   ingredient.UpdatedAt,
		Name:        ingredient.Name,
		Description: stringPointerFromSqlNullString(ingredient.Description),
		RecipeID:    ingredient.RecipeID,
	}
}

func (c *config) handleGetIngredients() http.HandlerFunc {
	type response []ingredientResponse

	return func(w http.ResponseWriter, r *http.Request) {
		var user database.User

		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		idString := chi.URLParam(r, "recipe_id")

		recipeID, err := strconv.Atoi(idString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Recipe id is not an integer")
			return
		}

		recipe, err := c.DB.GetRecipe(r.Context(), int64(recipeID))

		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Recipe does not exist")
			return
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not get recipe from database")
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

		resBody := make(response, len(ingredients))

		for i, dbIngredient := range ingredients {
			ingredient := databaseIngredientToReponse(dbIngredient)
			resBody[i] = ingredient
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}

}

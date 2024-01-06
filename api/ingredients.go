package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type ingredientResponse struct {
	ID          int64           `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Name        string          `json:"name"`
	Measure     measureResponse `json:"measure"`
	Description string          `json:"description,omitempty"`
	RecipeID    int64           `json:"recipe_id"`
}

type measureResponse struct {
	OriginalAmount float64 `json:"amount"`
	OriginalUnits  string  `json:"units"`
	StandardAmount float64 `json:"standard_amount"`
	StandardUnits  string  `json:"standard_units"`
}

func databaseIngredientToReponse(ingredient database.Ingredient) ingredientResponse {
	return ingredientResponse{
		ID:        ingredient.ID,
		CreatedAt: ingredient.CreatedAt,
		UpdatedAt: ingredient.UpdatedAt,
		Name:      ingredient.Name,
		Measure: measureResponse{
			OriginalAmount: ingredient.Amount,
			OriginalUnits:  ingredient.Units,
			StandardAmount: ingredient.StandardAmount,
			StandardUnits:  ingredient.StandardUnits,
		},
		Description: ingredient.Description.String,
		RecipeID:    ingredient.RecipeID,
	}
}

func (c *Config) handleGetIngredients() http.HandlerFunc {
	type response []ingredientResponse

	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(ContextUserKey).(database.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		idString := chi.URLParam(r, "recipe_id")

		recipeID, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Recipe id is not an integer")
			return
		}

		recipe, err := c.Domain.GetRecipe(r.Context(), user, recipeID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		ingredients, err := c.Domain.GetIngredientsForRecipe(r.Context(), user, recipe)
		if err != nil {
			respondWithDomainError(w, err)
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

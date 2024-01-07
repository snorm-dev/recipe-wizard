package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/domain"
)

type recipeResponse struct {
	ID          int64                `json:"id"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Url         string               `json:"url,omitempty"`
	PrepTime    string               `json:"prep_time,omitempty"`
	CookTime    string               `json:"cook_time,omitempty"`
	TotalTime   string               `json:"total_time,omitempty"`
	OwnerId     int64                `json:"owner_id"`
	Ingredients []ingredientResponse `json:"ingredients,omitempty"`
}

func domainRecipeToResponse(recipe domain.Recipe, ingredients []domain.Ingredient) recipeResponse {
	var responseIngredients []ingredientResponse

	if ingredients != nil {
		responseIngredients = make([]ingredientResponse, 0, len(ingredients))
		for _, ingredient := range ingredients {
			responseIngredients = append(responseIngredients, domainIngredientToReponse(ingredient))
		}
	}

	return recipeResponse{
		ID:          recipe.ID,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   recipe.UpdatedAt,
		Name:        recipe.Name,
		Description: (recipe.Description),
		Url:         (recipe.Url),
		PrepTime:    (recipe.PrepTime),
		CookTime:    (recipe.CookTime),
		TotalTime:   (recipe.TotalTime),
		OwnerId:     recipe.OwnerID,
		Ingredients: responseIngredients,
	}
}

func (c *Config) handlePostRecipe() http.HandlerFunc {
	type request struct {
		Url string `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		reqBody := request{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Unable to parse json body")
			return
		}

		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		recipe, err := c.Domain.CreateRecipeFromUrl(r.Context(), user, reqBody.Url)
		if err != nil {
			respondWithDomainError(w, err)
		}

		ingredients, err := c.Domain.GetIngredientsForRecipe(r.Context(), user, recipe)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		resBody := domainRecipeToResponse(recipe, ingredients)

		respondWithJSON(w, http.StatusCreated, &resBody)
	}
}

func (c *Config) handleGetRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "recipe_id")

		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		recipe, err := c.Domain.GetRecipe(r.Context(), user, id)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		var ingredients []domain.Ingredient
		if r.URL.Query().Has("return-ingredients") {
			ingredients, err = c.Domain.GetIngredientsForRecipe(r.Context(), user, recipe)
			if err != nil {
				respondWithDomainError(w, err)
				return
			}
		}

		resBody := domainRecipeToResponse(recipe, ingredients)

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

func (c *Config) handleGetRecipes() http.HandlerFunc {
	type response []recipeResponse

	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		recipes, err := c.Domain.GetRecipesForUser(r.Context(), user)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		resBody := make(response, len(recipes))

		for i, recipe := range recipes {

			var ingredients []domain.Ingredient
			if r.URL.Query().Has("return-ingredients") {
				ingredients, err = c.Domain.GetIngredientsForRecipe(r.Context(), user, recipe)
				if err != nil {
					respondWithDomainError(w, err)
					return
				}
			}

			r := domainRecipeToResponse(recipe, ingredients)
			resBody[i] = r
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

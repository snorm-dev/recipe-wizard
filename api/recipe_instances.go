package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/domain"
)

type recipeInstanceResponse struct {
	ID            int64          `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	GroceryListID int64          `json:"grocery_list_id"`
	RecipeID      int64          `json:"recipe_id"`
	Items         []itemResponse `json:"items"`
}

func domainRecipeInstanceToResponse(ri domain.RecipeInstance, its []domain.Item) recipeInstanceResponse {
	items := make([]itemResponse, len(its))
	for idx, it := range its {
		items[idx] = domainItemToResponse(it)
	}
	return recipeInstanceResponse{
		ID:            ri.ID,
		CreatedAt:     ri.CreatedAt,
		UpdatedAt:     ri.UpdatedAt,
		GroceryListID: ri.GroceryListID,
		RecipeID:      ri.Recipe.ID,
		Items:         items,
	}
}

func (c *Config) handlePostRecipeInGroceryList() http.HandlerFunc {

	type request struct {
		RecipeID int64 `json:"recipe_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		idString := chi.URLParam(r, "grocery_list_id")

		glID, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		groceryList, err := c.Domain.GetGroceryList(r.Context(), user, glID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		reqBody := request{}

		err = json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		recipeInstance, err := c.Domain.CreateRecipeInstance(r.Context(), user, groceryList, reqBody.RecipeID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		items, err := c.Domain.GetItemsForRecipeInstance(r.Context(), recipeInstance)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		var resBody = domainRecipeInstanceToResponse(recipeInstance, items)

		respondWithJSON(w, http.StatusCreated, &resBody)
	}
}

func (c *Config) handleGetRecipesInGroceryList() http.HandlerFunc {

	type response = []recipeInstanceResponse

	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		idString := chi.URLParam(r, "grocery_list_id")

		glID, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		groceryList, err := c.Domain.GetGroceryList(r.Context(), user, glID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		recipeInstances, err := c.Domain.GetRecipeInstancesInGroceryList(r.Context(), groceryList)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		var resBody response = make([]recipeInstanceResponse, len(recipeInstances))

		for i, recipeInstance := range recipeInstances {

			items, err := c.Domain.GetItemsForRecipeInstance(r.Context(), recipeInstance)
			if err != nil {
				respondWithDomainError(w, err)
				return
			}

			r := domainRecipeInstanceToResponse(recipeInstance, items)
			resBody[i] = r
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

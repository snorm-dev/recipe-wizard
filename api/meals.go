package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/domain"
)

type mealResponse struct {
	ID            int64          `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	GroceryListID int64          `json:"grocery_list_id"`
	RecipeID      int64          `json:"recipe_id"`
	Items         []itemResponse `json:"items"`
}

func domainMealToResponse(m domain.Meal, its []domain.Item) mealResponse {
	items := make([]itemResponse, len(its))
	for idx, it := range its {
		items[idx] = domainItemToResponse(it)
	}
	return mealResponse{
		ID:            m.ID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		GroceryListID: m.GroceryListID,
		RecipeID:      m.Recipe.ID,
		Items:         items,
	}
}

func (c *Config) handlePostMealInGroceryList() http.HandlerFunc {

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

		meal, err := c.Domain.CreateMeal(r.Context(), user, groceryList, reqBody.RecipeID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		items, err := c.Domain.GetItemsForMeal(r.Context(), meal)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		var resBody = domainMealToResponse(meal, items)

		respondWithJSON(w, http.StatusCreated, &resBody)
	}
}

func (c *Config) handleGetMealsInGroceryList() http.HandlerFunc {

	type response = []mealResponse

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

		meals, err := c.Domain.GetMealsInGroceryList(r.Context(), groceryList)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		var resBody response = make([]mealResponse, len(meals))

		for i, meal := range meals {

			items, err := c.Domain.GetItemsForMeal(r.Context(), meal)
			if err != nil {
				respondWithDomainError(w, err)
				return
			}

			m := domainMealToResponse(meal, items)
			resBody[i] = m
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

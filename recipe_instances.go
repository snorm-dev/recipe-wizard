package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type recipeInstanceResponse struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	GroceryListID int64     `json:"grocery_list_id"`
	RecipeID      int64     `json:"recipe_id"`
}

func databaseRecipeInstanceToResponse(ri database.RecipeInstance) recipeInstanceResponse {
	return recipeInstanceResponse{
		ID:            ri.ID,
		CreatedAt:     ri.CreatedAt,
		UpdatedAt:     ri.UpdatedAt,
		GroceryListID: ri.GroceryListID,
		RecipeID:      ri.RecipeID,
	}
}

func (c *config) handlePostRecipeInGroceryList() http.HandlerFunc {

	type request struct {
		RecipeID int64 `json:"recipe_id"`
	}

	type response = recipeInstanceResponse

	return func(w http.ResponseWriter, r *http.Request) {
		var user database.User
		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		idString := chi.URLParam(r, "grocery_list_id")

		glID, err := strconv.Atoi(idString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		groceryList, err := c.DB.GetGroceryList(r.Context(), int64(glID))

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if user.ID != groceryList.OwnerID {
			respondWithError(w, http.StatusForbidden, "User does not own grocery list.")
			return
		}

		reqBody := request{}

		err = json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		now := time.Now()

		result, err := c.DB.AddRecipeToGroceryList(r.Context(), database.AddRecipeToGroceryListParams{
			CreatedAt:     now,
			UpdatedAt:     now,
			RecipeID:      reqBody.RecipeID,
			GroceryListID: groceryList.ID,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create recipe instance in db: %v", err))
			return
		}

		id, err := result.LastInsertId()

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not retrieve id: %v", err))
		}

		recipeInstance, err := c.DB.GetRecipeInstance(r.Context(), id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get recipe instance from db: %v", err))
			return
		}

		var resBody response = databaseRecipeInstanceToResponse(recipeInstance)

		respondWithJSON(w, http.StatusCreated, &resBody)
	}
}

func (c *config) handleGetRecipesInGroceryList() http.HandlerFunc {

	type response = []recipeResponse

	return func(w http.ResponseWriter, r *http.Request) {

		var user database.User
		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		idString := chi.URLParam(r, "grocery_list_id")

		glID, err := strconv.Atoi(idString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		groceryList, err := c.DB.GetGroceryList(r.Context(), int64(glID))

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if user.ID != groceryList.OwnerID {
			respondWithError(w, http.StatusForbidden, "User does not own grocery list.")
			return
		}

		recipes, err := c.DB.GetRecipesInGroceryList(r.Context(), groceryList.ID)

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

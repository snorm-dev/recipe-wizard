package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type recipeInstanceResponse struct {
	ID                  int64                        `json:"id"`
	CreatedAt           time.Time                    `json:"created_at"`
	UpdatedAt           time.Time                    `json:"updated_at"`
	GroceryListID       int64                        `json:"grocery_list_id"`
	RecipeID            int64                        `json:"recipe_id"`
	IngredientInstances []ingredientInstanceResponse `json:"ingredient_instances"`
}

func databaseRecipeInstanceToResponse(ri database.RecipeInstance, iis []database.IngredientInstance) recipeInstanceResponse {
	instances := make([]ingredientInstanceResponse, len(iis))
	for _, ii := range iis {
		instances = append(instances, databaseIngredientInstanceToResponse(ii))
	}
	return recipeInstanceResponse{
		ID:                  ri.ID,
		CreatedAt:           ri.CreatedAt,
		UpdatedAt:           ri.UpdatedAt,
		GroceryListID:       ri.GroceryListID,
		RecipeID:            ri.RecipeID,
		IngredientInstances: instances,
	}
}

func (c *config) handlePostRecipeInGroceryList() http.HandlerFunc {

	type request struct {
		RecipeID int64 `json:"recipe_id"`
	}

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

		result, err := c.DB.CreateRecipeInstance(r.Context(), database.CreateRecipeInstanceParams{
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

		ingredients, err := c.DB.GetIngredientsForRecipe(r.Context(), reqBody.RecipeID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		for _, ingredient := range ingredients {
			now := time.Now()
			_, err := c.DB.CreateIngredientInstance(r.Context(), database.CreateIngredientInstanceParams{
				CreatedAt:        now,
				UpdatedAt:        now,
				IngredientID:     sql.NullInt64{Int64: ingredient.ID, Valid: true},
				GroceryListID:    groceryList.ID,
				RecipeInstanceID: sql.NullInt64{Int64: recipeInstance.ID, Valid: true},
			})

			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		ingredientInstances, err := c.DB.GetIngredientInstancesForRecipeInstance(r.Context(), sql.NullInt64{Valid: true, Int64: recipeInstance.ID})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var resBody = databaseRecipeInstanceToResponse(recipeInstance, ingredientInstances)

		respondWithJSON(w, http.StatusCreated, &resBody)
	}
}

func (c *config) handleGetRecipesInGroceryList() http.HandlerFunc {

	type response = []recipeInstanceResponse

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

		recipeInstances, err := c.DB.GetRecipeInstancesInGroceryList(r.Context(), groceryList.ID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var resBody response = make([]recipeInstanceResponse, len(recipeInstances))

		for i, recipe := range recipeInstances {

			ingredients, err := c.DB.GetIngredientInstancesForRecipeInstance(r.Context(), sql.NullInt64{Valid: true, Int64: recipe.ID})

			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			r := databaseRecipeInstanceToResponse(recipe, ingredients)
			resBody[i] = r
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

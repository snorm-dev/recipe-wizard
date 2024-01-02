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

type groceryListResponse struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	OwnerID   int64     `json:"owner_id"`
}

func databaseGroceryListToResponse(gl database.GroceryList) groceryListResponse {
	return groceryListResponse{
		ID:        gl.ID,
		CreatedAt: gl.CreatedAt,
		UpdatedAt: gl.UpdatedAt,
		Name:      gl.Name,
		OwnerID:   gl.OwnerID,
	}
}

func (c *config) handlePostGroceryList() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := request{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		now := time.Now()

		var user database.User
		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		result, err := c.DB.CreateGroceryList(r.Context(), database.CreateGroceryListParams{
			CreatedAt: now,
			UpdatedAt: now,
			Name:      reqBody.Name,
			OwnerID:   user.ID,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create grocery list in db: %v", err))
			return
		}

		id, err := result.LastInsertId()

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not retrieve id: %v", err))
		}

		groceryList, err := c.DB.GetGroceryList(r.Context(), id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get grocery list from db: %v", err))
			return
		}

		respondWithJSON(w, http.StatusCreated, &groceryList)
	}
}

func (c *config) handleGetGroceryList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "grocery_list_id")

		id, err := strconv.Atoi(idString)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		groceryList, err := c.DB.GetGroceryList(r.Context(), int64(id))

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

		if user.ID != groceryList.OwnerID {
			respondWithError(w, http.StatusForbidden, "User does not own grocery list.")
			return
		}

		resBody := databaseGroceryListToResponse(groceryList)

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

func (c *config) handleGetGroceryLists() http.HandlerFunc {
	type response []groceryListResponse

	return func(w http.ResponseWriter, r *http.Request) {

		var user database.User
		if value := r.Context().Value(ContextUserKey); value != nil {
			user = value.(database.User)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Could not get user")
			return
		}

		groceryLists, err := c.DB.GetGroceryListsForUser(r.Context(), user.ID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := make(response, len(groceryLists))

		for i, dbGroceryList := range groceryLists {
			r := databaseGroceryListToResponse(dbGroceryList)
			resBody[i] = r
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

func (c *config) handleGetIngredientsInGroceryList() http.HandlerFunc {
	type response []ingredientResponse

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

		ingredients, err := c.DB.GetIngredientsInGroceryList(r.Context(), groceryList.ID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		resBody := make(response, len(ingredients))

		for i, ingredient := range ingredients {
			resIngredient := databaseIngredientToReponse(ingredient)
			resBody[i] = resIngredient
		}

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

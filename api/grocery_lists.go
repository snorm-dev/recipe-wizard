package api

import (
	"encoding/json"
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

func (c *Config) handlePostGroceryList() http.HandlerFunc {
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

		user, ok := r.Context().Value(ContextUserKey).(database.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		groceryList, err := c.Domain.CreateGroceryList(r.Context(), user, reqBody.Name)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		respondWithJSON(w, http.StatusCreated, &groceryList)
	}
}

func (c *Config) handleGetGroceryList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "grocery_list_id")

		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		user, ok := r.Context().Value(ContextUserKey).(database.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		groceryList, err := c.Domain.GetGroceryList(r.Context(), user, id)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		resBody := databaseGroceryListToResponse(groceryList)

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

func (c *Config) handleGetGroceryLists() http.HandlerFunc {
	type response []groceryListResponse

	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(ContextUserKey).(database.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		groceryLists, err := c.Domain.GetGroceryListsForUser(r.Context(), user)
		if err != nil {
			respondWithDomainError(w, err)
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

/*func (c *Config) handleGetIngredientsInGroceryList() http.HandlerFunc {
	type response []ingredientResponse

	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(database.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
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
}*/

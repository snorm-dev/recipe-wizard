package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/snorman7384/recipe-wizard/internal/database"
)

func (c *config) handlePostGroceryList() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	type response struct {
		Name    string `json:"name"`
		OwnerID int64  `json:"owner_id"`
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

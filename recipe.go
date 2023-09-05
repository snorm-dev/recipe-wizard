package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kkyr/go-recipe/pkg/recipe"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

func (c *config) handlePostRecipe() http.HandlerFunc {
	type request struct {
		Url string `json:"url"`
	}
	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Name        string    `json:"name"`
		Description *string   `json:"description,omitempty"`
		Url         *string   `json:"url,omitempty"`
		PrepTime    *string   `json:"prep_time,omitempty"`
		CookTime    *string   `json:"cook_time,omitempty"`
		TotalTime   *string   `json:"total_time,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := request{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// get recipe data
		s, err := recipe.ScrapeURL(reqBody.Url)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse recipe from url: %v", err))
		}

		name, ok := s.Name()
		if !ok {
			name = reqBody.Url
		}

		str, ok := s.Description()
		description := sql.NullString{
			String: str,
			Valid:  ok,
		}

		t, ok := s.PrepTime()
		prepTime := sql.NullString{
			String: t.String(),
			Valid:  ok,
		}

		t, ok = s.CookTime()
		cookTime := sql.NullString{
			String: t.String(),
			Valid:  ok,
		}

		t, ok = s.TotalTime()
		totalTime := sql.NullString{
			String: t.String(),
			Valid:  ok,
		}

		id := uuid.New().String()
		now := time.Now()
		url := sql.NullString{
			String: reqBody.Url,
			Valid:  true,
		}

		err = c.DB.CreateRecipe(r.Context(), database.CreateRecipeParams{
			ID:          id,
			CreatedAt:   now,
			UpdatedAt:   now,
			Url:         url,
			Name:        name,
			Description: description,
			CookTime:    cookTime,
			PrepTime:    prepTime,
			TotalTime:   totalTime,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create recipe in db: %v", err))
			return
		}

		recipe, err := c.DB.GetRecipe(r.Context(), id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not get recipe from db: %v", err))
			return
		}

		uuid, err := uuid.Parse(recipe.ID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not parse uuid: %v", err))
			return
		}

		resBody := response{
			ID:        uuid,
			CreatedAt: recipe.CreatedAt,
			UpdatedAt: recipe.UpdatedAt,
			Name:      name,
		}

		if recipe.Description.Valid {
			resBody.Description = &recipe.Description.String
		}

		if recipe.Url.Valid {
			resBody.Url = &recipe.Url.String
		}

		if recipe.PrepTime.Valid {
			resBody.PrepTime = &recipe.PrepTime.String
		}

		if recipe.CookTime.Valid {
			resBody.CookTime = &recipe.CookTime.String
		}

		if recipe.TotalTime.Valid {
			resBody.TotalTime = &recipe.TotalTime.String
		}

		err = respondWithJSON(w, http.StatusCreated, &resBody)

		if err != nil {
			log.Println("Error responding to request: ", err)
		}
	}
}

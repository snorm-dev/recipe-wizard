package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/domain"
)

type ingredientInstanceResponse struct {
	ID               int64              `json:"id"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	GroceryListID    int64              `json:"grocery_list_id"`
	RecipeInstanceID int64              `json:"recipe_instance_id,omitempty"` // 0 is never a sql id, so we can treat 0 as "no recipe instance"
	IngredientData   ingredientResponse `json:"ingredient_data"`
}

type ingredientGroupResponse struct {
	Name      string                       `json:"name"`
	Total     float64                      `json:"total"`
	Units     string                       `json:"units"`
	Instances []ingredientInstanceResponse `json:"ingredient_instances"`
}

func domainIngredientInstanceToResponse(ii domain.IngredientInstance) ingredientInstanceResponse {
	return ingredientInstanceResponse{
		ID:               ii.ID,
		CreatedAt:        ii.CreatedAt,
		UpdatedAt:        ii.UpdatedAt,
		GroceryListID:    ii.GroceryListID,
		RecipeInstanceID: ii.RecipeInstanceID,
		IngredientData:   domainIngredientToReponse(ii.Ingredient),
	}
}

func domainIngredientGroupToResponse(ig domain.IngredientGroup) ingredientGroupResponse {
	instances := make([]ingredientInstanceResponse, len(ig.Instances))
	for i, ii := range ig.Instances {
		instances[i] = domainIngredientInstanceToResponse(ii)
	}
	return ingredientGroupResponse{
		Name:      ig.Name,
		Total:     ig.Total,
		Units:     ig.Units.String(),
		Instances: instances,
	}
}

func (c *Config) handleGetIngredientInstancesForGroceryList() http.HandlerFunc {
	type response struct {
		Instances []ingredientInstanceResponse `json:"ingredient_instances,omitempty"`
		Groups    []ingredientGroupResponse    `json:"ingredient_groups,omitempty"`
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

		if r.URL.Query().Has("grouped") && r.URL.Query().Has("ungrouped") {
			respondWithError(w, http.StatusConflict, "Conflicting query parameters")
			return
		}

		resBody := response{}
		if r.URL.Query().Has("grouped") {
			ingredientGroups, err := c.Domain.GetIngredientGroupsForGroceryList(r.Context(), groceryList)
			if err != nil {
				respondWithDomainError(w, err)
				return
			}

			for _, ingredientGroup := range ingredientGroups {
				r := domainIngredientGroupToResponse(ingredientGroup)
				resBody.Groups = append(resBody.Groups, r)
			}
		} else {
			ingredientInstances, err := c.Domain.GetIngredientInstancesForGroceryList(r.Context(), groceryList)
			if err != nil {
				respondWithDomainError(w, err)
				return
			}

			for _, ingredientInstance := range ingredientInstances {
				r := domainIngredientInstanceToResponse(ingredientInstance)
				resBody.Instances = append(resBody.Instances, r)
			}
		}
		respondWithJSON(w, http.StatusOK, resBody)
	}
}

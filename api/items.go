package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/domain"
)

type itemResponse struct {
	ID               int64              `json:"id"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	GroceryListID    int64              `json:"grocery_list_id"`
	RecipeInstanceID int64              `json:"recipe_instance_id,omitempty"` // 0 is never a sql id, so we can treat 0 as "no recipe instance"
	IngredientData   ingredientResponse `json:"ingredient_data"`
}

type itemGroupResponse struct {
	Name  string         `json:"name"`
	Total float64        `json:"total"`
	Units string         `json:"units"`
	Items []itemResponse `json:"items"`
}

func domainItemToResponse(ii domain.Item) itemResponse {
	return itemResponse{
		ID:               ii.ID,
		CreatedAt:        ii.CreatedAt,
		UpdatedAt:        ii.UpdatedAt,
		GroceryListID:    ii.GroceryListID,
		RecipeInstanceID: ii.RecipeInstanceID,
		IngredientData:   domainIngredientToReponse(ii.Ingredient),
	}
}

func domainItemGroupToResponse(ig domain.ItemGroup) itemGroupResponse {
	items := make([]itemResponse, len(ig.Items))
	for i, ii := range ig.Items {
		items[i] = domainItemToResponse(ii)
	}
	return itemGroupResponse{
		Name:  ig.Name,
		Total: ig.Total,
		Units: ig.Units.String(),
		Items: items,
	}
}

func (c *Config) handleGetItemsForGroceryList() http.HandlerFunc {
	type response struct {
		Items  []itemResponse      `json:"items,omitempty"`
		Groups []itemGroupResponse `json:"item_groups,omitempty"`
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
			itemGroups, err := c.Domain.GetItemGroupsForGroceryList(r.Context(), groceryList)
			if err != nil {
				respondWithDomainError(w, err)
				return
			}

			for _, itemGroup := range itemGroups {
				r := domainItemGroupToResponse(itemGroup)
				resBody.Groups = append(resBody.Groups, r)
			}
		} else {
			items, err := c.Domain.GetItemsForGroceryList(r.Context(), groceryList)
			if err != nil {
				respondWithDomainError(w, err)
				return
			}

			for _, item := range items {
				r := domainItemToResponse(item)
				resBody.Items = append(resBody.Items, r)
			}
		}
		respondWithJSON(w, http.StatusOK, resBody)
	}
}

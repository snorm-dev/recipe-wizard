package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/snorman7384/recipe-wizard/domain"
	"github.com/snorman7384/recipe-wizard/ingparse"
)

type itemResponse struct {
	ID            int64           `json:"id"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	GroceryListID int64           `json:"grocery_list_id"`
	MealID        int64           `json:"meal_id,omitempty"` // 0 is never a sql id, so we can treat 0 as "no meal"
	IngredientID  int64           `json:"ingredient_id,omitempty"`
	Name          string          `json:"name"`
	Description   string          `json:"description,omitempty"`
	Measure       measureResponse `json:"measure"`
	Status        string          `json:"status"`
}

type itemGroupResponse struct {
	Name   string                            `json:"name"`
	Totals map[ingparse.StandardUnit]float64 `json:"totals,omitempty"`
	Items  []itemResponse                    `json:"items,omitempty"`
}

func domainItemToResponse(it domain.Item) itemResponse {
	return itemResponse{
		ID:            it.ID,
		CreatedAt:     it.CreatedAt,
		UpdatedAt:     it.UpdatedAt,
		GroceryListID: it.GroceryListID,
		MealID:        it.MealID,
		IngredientID:  it.IngredientID,
		Name:          it.Name,
		Description:   it.Description,
		Measure: measureResponse{
			OriginalAmount: it.Amount,
			OriginalUnits:  it.Units,
			StandardAmount: it.StandardAmount,
			StandardUnits:  it.StandardUnits.String(),
		},
		Status: it.Status.String(),
	}
}

func domainItemGroupToResponse(ig domain.ItemGroup) itemGroupResponse {
	items := make([]itemResponse, len(ig.Items))
	for i, it := range ig.Items {
		items[i] = domainItemToResponse(it)
	}
	return itemGroupResponse{
		Name:   ig.Name,
		Totals: ig.Totals,
		Items:  items,
	}
}

func (c *Config) handlePostItem() http.HandlerFunc {
	type request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Units       string  `json:"units"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := request{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

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

		item, err := c.Domain.CreateItem(r.Context(), groceryList, reqBody.Name, reqBody.Description, float64(reqBody.Amount), reqBody.Units)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		resBody := domainItemToResponse(item)

		respondWithJSON(w, http.StatusOK, resBody)
	}
}

func (c *Config) handleGetItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		idString := chi.URLParam(r, "item_id")

		itemID, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		item, err := c.Domain.GetItem(r.Context(), user, itemID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		respondWithJSON(w, http.StatusOK, domainItemToResponse(item))
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

func (c *Config) handleGetItemsForGroceryListByName() http.HandlerFunc {
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

		name := chi.URLParam(r, "item_name")

		groceryList, err := c.Domain.GetGroceryList(r.Context(), user, glID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		itemGroup, err := c.Domain.GetItemGroupForGroceryListByName(r.Context(), groceryList, name)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		respondWithJSON(w, http.StatusOK, domainItemGroupToResponse(itemGroup))
	}
}

func (c *Config) handleMarkItemStatus() http.HandlerFunc {
	type request struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		reqBody := request{}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		status, err := domain.ItemStatusFromString(reqBody.Status)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		idString := chi.URLParam(r, "item_id")

		itemID, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		item, err := c.Domain.GetItem(r.Context(), user, itemID)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		item, err = c.Domain.MarkItemStatus(r.Context(), item, status)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		respondWithJSON(w, http.StatusOK, domainItemToResponse(item))
	}
}

package domain

import (
	"bytes"
	"errors"
	"time"

	"github.com/snorman7384/recipe-wizard/ingparse"
)

type GroceryList struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	OwnerID   int64
}

type Ingredient struct {
	ID             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Description    string
	RecipeID       int64
	Amount         float64
	Units          string
	StandardAmount float64
	StandardUnits  ingparse.StandardUnit
}

type Item struct {
	ID             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	GroceryListID  int64
	MealID         int64
	IngredientID   int64
	Name           string
	Description    string
	Amount         float64
	Units          string
	StandardAmount float64
	StandardUnits  ingparse.StandardUnit
	Status         ItemStatus
}

type Recipe struct {
	ID          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	Url         string
	PrepTime    string
	CookTime    string
	TotalTime   string
	OwnerID     int64
}

type Meal struct {
	ID            int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	GroceryListID int64
	Recipe        Recipe
}

type User struct {
	ID             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Username       string
	HashedPassword string
	FirstName      string
	LastName       string
}

type ItemGroup struct {
	Name   string
	Totals map[ingparse.StandardUnit]float64
	Items  []Item
}

type ItemStatus int

const (
	_ ItemStatus = iota
	Incomplete
	Complete
)

func (s ItemStatus) String() string {
	if s == Incomplete {
		return "incomplete"
	}
	if s == Complete {
		return "complete"
	}
	return "<error>"
}

func (s ItemStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s ItemStatus) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func ItemStatusFromString(s string) (ItemStatus, error) {
	if s == "incomplete" {
		return Incomplete, nil
	}
	if s == "complete" {
		return Complete, nil
	}
	return 0, errors.New("invalid status string")
}

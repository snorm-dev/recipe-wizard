package domain

import (
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
	ID               int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	GroceryListID    int64
	RecipeInstanceID int64
	Ingredient       Ingredient
	Name             string
	Description      string
	Amount           float64
	Units            string
	StandardAmount   float64
	StandardUnits    ingparse.StandardUnit
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

type RecipeInstance struct {
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
	Name  string
	Total float64
	Units ingparse.StandardUnit
	Items []Item
}

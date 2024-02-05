// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"database/sql"
	"time"
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
	Description    sql.NullString
	RecipeID       int64
	Amount         float64
	Units          string
	StandardAmount float64
	StandardUnits  string
}

type Item struct {
	ID               int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	GroceryListID    int64
	RecipeInstanceID sql.NullInt64
	IngredientID     sql.NullInt64
	Name             string
	Description      sql.NullString
	Amount           float64
	Units            string
	StandardAmount   float64
	StandardUnits    string
}

type Recipe struct {
	ID          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description sql.NullString
	Url         sql.NullString
	PrepTime    sql.NullString
	CookTime    sql.NullString
	TotalTime   sql.NullString
	OwnerID     int64
}

type RecipeInstance struct {
	ID            int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	GroceryListID int64
	RecipeID      int64
}

type User struct {
	ID             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Username       string
	HashedPassword string
	FirstName      sql.NullString
	LastName       sql.NullString
}

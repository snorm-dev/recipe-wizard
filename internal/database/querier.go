// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateGroceryList(ctx context.Context, arg CreateGroceryListParams) (GroceryList, error)
	CreateIngredient(ctx context.Context, arg CreateIngredientParams) (Ingredient, error)
	CreateItem(ctx context.Context, arg CreateItemParams) (Item, error)
	CreateMeal(ctx context.Context, arg CreateMealParams) (Meal, error)
	CreateRecipe(ctx context.Context, arg CreateRecipeParams) (Recipe, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetExtendedItem(ctx context.Context, id int64) (GetExtendedItemRow, error)
	GetExtendedItemsForGroceryList(ctx context.Context, groceryListID int64) ([]GetExtendedItemsForGroceryListRow, error)
	GetExtendedItemsForMeal(ctx context.Context, mealID sql.NullInt64) ([]GetExtendedItemsForMealRow, error)
	GetExtendedMeal(ctx context.Context, id int64) (GetExtendedMealRow, error)
	GetExtendedMealsInGroceryList(ctx context.Context, groceryListID int64) ([]GetExtendedMealsInGroceryListRow, error)
	GetGroceryList(ctx context.Context, id int64) (GroceryList, error)
	GetGroceryListsForUser(ctx context.Context, ownerID int64) ([]GroceryList, error)
	GetIngredient(ctx context.Context, id int64) (Ingredient, error)
	GetIngredientsForRecipe(ctx context.Context, recipeID int64) ([]Ingredient, error)
	GetItem(ctx context.Context, id int64) (Item, error)
	GetItemAndGroceryList(ctx context.Context, id int64) (GetItemAndGroceryListRow, error)
	GetItemsForGroceryList(ctx context.Context, groceryListID int64) ([]Item, error)
	GetItemsForGroceryListByName(ctx context.Context, arg GetItemsForGroceryListByNameParams) ([]Item, error)
	GetItemsForMeal(ctx context.Context, mealID sql.NullInt64) ([]Item, error)
	GetMeal(ctx context.Context, id int64) (Meal, error)
	GetMealsInGroceryList(ctx context.Context, groceryListID int64) ([]Meal, error)
	GetRecipe(ctx context.Context, id int64) (Recipe, error)
	GetRecipesForUser(ctx context.Context, ownerID int64) ([]Recipe, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	SetIsComplete(ctx context.Context, arg SetIsCompleteParams) error
}

var _ Querier = (*Queries)(nil)

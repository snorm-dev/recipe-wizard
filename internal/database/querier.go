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
	CreateRecipe(ctx context.Context, arg CreateRecipeParams) (Recipe, error)
	CreateRecipeInstance(ctx context.Context, arg CreateRecipeInstanceParams) (RecipeInstance, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetExtendedItem(ctx context.Context, id int64) (GetExtendedItemRow, error)
	GetExtendedItemsForGroceryList(ctx context.Context, groceryListID int64) ([]GetExtendedItemsForGroceryListRow, error)
	GetExtendedItemsForRecipeInstance(ctx context.Context, recipeInstanceID sql.NullInt64) ([]GetExtendedItemsForRecipeInstanceRow, error)
	GetExtendedRecipeInstance(ctx context.Context, id int64) (GetExtendedRecipeInstanceRow, error)
	GetExtendedRecipeInstancesInGroceryList(ctx context.Context, groceryListID int64) ([]GetExtendedRecipeInstancesInGroceryListRow, error)
	GetGroceryList(ctx context.Context, id int64) (GroceryList, error)
	GetGroceryListsForUser(ctx context.Context, ownerID int64) ([]GroceryList, error)
	GetIngredient(ctx context.Context, id int64) (Ingredient, error)
	GetIngredientsForRecipe(ctx context.Context, recipeID int64) ([]Ingredient, error)
	GetItem(ctx context.Context, id int64) (Item, error)
	GetItemAndGroceryList(ctx context.Context, id int64) (GetItemAndGroceryListRow, error)
	GetItemsForGroceryList(ctx context.Context, groceryListID int64) ([]Item, error)
	GetItemsForGroceryListByName(ctx context.Context, arg GetItemsForGroceryListByNameParams) ([]Item, error)
	GetItemsForRecipeInstance(ctx context.Context, recipeInstanceID sql.NullInt64) ([]Item, error)
	GetRecipe(ctx context.Context, id int64) (Recipe, error)
	GetRecipeInstance(ctx context.Context, id int64) (RecipeInstance, error)
	GetRecipeInstancesInGroceryList(ctx context.Context, groceryListID int64) ([]RecipeInstance, error)
	GetRecipesForUser(ctx context.Context, ownerID int64) ([]Recipe, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	SetIsComplete(ctx context.Context, arg SetIsCompleteParams) error
}

var _ Querier = (*Queries)(nil)

package domain

import (
	"database/sql"

	"github.com/snorman7384/recipe-wizard/ingparse"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type Config struct {
	DB               *sql.DB
	IngredientParser ingparse.IngredientParser
}

func (c *Config) Querier() *database.Queries {
	return database.New(c.DB)
}

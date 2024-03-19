package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/snorman7384/recipe-wizard/api"
	"github.com/snorman7384/recipe-wizard/domain"
	"github.com/snorman7384/recipe-wizard/ingparse"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Could not load .env file: ", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Could not load custom port.")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		log.Fatal("Could not load database url")
	}

	db, err := sql.Open("libsql", dbUrl)

	if err != nil {
		log.Fatal("Could not open database connection")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("Could not locate JWT secret")
	}

	c := api.Config{
		Domain: domain.Config{
			DB:               db,
			IngredientParser: ingparse.SchollzParser{},
		},
		JwtSecret: []byte(jwtSecret),
		Port:      port,
	}

	c.Serve()
}

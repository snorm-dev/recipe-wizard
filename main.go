package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type config struct {
	DB *database.Queries
}

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

	db, err := sql.Open("mysql", dbUrl)

	if err != nil {
		log.Fatal("Could not open database connection")
	}

	c := config{
		DB: database.New(db),
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()
	r.Mount("/v1", v1)

	v1.Get("/ping", c.handlePing())
	v1.Post("/recipes", c.handlePostRecipe())
	v1.Get("/recipes", c.handleGetRecipes())
	v1.Get("/recipes/{recipe_id}", c.handleGetRecipe())

	server := &http.Server{
		Addr:              "0.0.0.0:" + port,
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("Listening on port: ", port)
	log.Fatal(server.ListenAndServe())
}

func (c *config) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong!"))
		if err != nil {
			log.Println("Could not respond to ping request: ", err)
		}
	}
}

func respondWithJSON(w http.ResponseWriter, code int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	bytes, err := json.Marshal(body)

	if err != nil {
		return err
	}

	_, err = w.Write(bytes)

	return err
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	type response struct {
		Error string `json:"error"`
	}

	err := respondWithJSON(w, code, response{Error: message})

	if err != nil {
		log.Println("Could not respond with error: ", err)
	}
}

func stringPointerFromSqlNullString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

func sqlNullStringFromString(s string, ok bool) sql.NullString {
	return sql.NullString{Valid: ok, String: s}
}

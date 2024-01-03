package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/snorman7384/recipe-wizard/ingparse"
	"github.com/snorman7384/recipe-wizard/internal/database"
)

type config struct {
	DB               database.Querier
	JwtSecret        []byte
	IngredientParser ingparse.IngredientParser
}

type ContextKey string

const ContextUserKey ContextKey = "user-key"

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

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("Could not locate JWT secret")
	}

	c := config{
		DB:               database.New(db),
		JwtSecret:        []byte(jwtSecret),
		IngredientParser: ingparse.SchollzParser{},
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

	r.Get("/", c.handleIndex())

	v1 := chi.NewRouter()
	r.Mount("/v1", middlewareLogRequest(v1))

	v1.Get("/ping", c.handlePing())

	v1.Post("/recipes", c.middlewareExtractUser(c.handlePostRecipe()))
	v1.Get("/recipes", c.middlewareExtractUser(c.handleGetRecipes()))
	v1.Get("/recipes/{recipe_id}", c.middlewareExtractUser(c.handleGetRecipe()))

	v1.Get("/recipes/{recipe_id}/ingredients", c.middlewareExtractUser(c.handleGetIngredients()))

	v1.Post("/grocery-lists", c.middlewareExtractUser(c.handlePostGroceryList()))
	v1.Get("/grocery-lists", c.middlewareExtractUser(c.handleGetGroceryLists()))
	v1.Get("/grocery-lists/{grocery_list_id}", c.middlewareExtractUser(c.handleGetGroceryList()))

	v1.Post("/grocery-lists/{grocery_list_id}/recipes", c.middlewareExtractUser(c.handlePostRecipeInGroceryList()))
	v1.Get("/grocery-lists/{grocery_list_id}/recipes", c.middlewareExtractUser(c.handleGetRecipesInGroceryList()))

	v1.Get("/grocery-lists/{grocery_list_id}/ingredients", c.middlewareExtractUser(c.handleGetIngredientsInGroceryList()))

	v1.Post("/users", c.handlePostUser())
	v1.Post("/login", c.handleLogin())

	server := &http.Server{
		Addr:              "0.0.0.0:" + port,
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("Listening on port: ", port)
	log.Fatal(server.ListenAndServe())
}

func (c *config) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myPage := `
		<!DOCTYPE html>
		<html>
		<body>

		<h1>Recipe Wizard</h1>
		<p>The API is available beginning with path /v1 </p>
		<p>Use the following endpoints:</p>
		<ul>
		<li>POST /v1/users</li>
		<li>POST /v1/login</li>
		<li>GET/POST /v1/recipes</li>
		<li>GET /v1/recipes{id}</li>
		<li>GET/POST /v1/recipes/{id}/ingredients</li>
		<li>GET/POST /v1/grocery-lists</li>
		<li>GET /v1/grocery-lists{id}</li>
		<li>GET/POST /v1/grocery-lists/recipes</li>
		<li>GET /v1/grocery-lists/ingredients</li>
		</ul>

		</body>
		</html>
		`
		_, err := w.Write([]byte(myPage))

		if err != nil {
			log.Println("Could not respond to index request: ", err)
		}
	}
}

func (c *config) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong!"))
		if err != nil {
			log.Println("Could not respond to ping request: ", err)
		}
	}
}

func respondWithJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	bytes, err := json.Marshal(body)

	if err != nil {
		log.Println("Could not marshal json: ", err)
		return
	}

	_, err = w.Write(bytes)

	if err != nil {
		log.Println("Could not write json to output: ", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	type response struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, response{Error: message})
}

func stringPointerFromSqlNullString(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

func int64PointerFromSqlNullInt64(i sql.NullInt64) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

func sqlNullStringFromOkString(s string, ok bool) sql.NullString {
	return sql.NullString{Valid: ok, String: s}
}

func sqlNullStringFromStringPointer(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	} else {
		return sql.NullString{Valid: true, String: *s}
	}
}

func middlewareLogRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: %v %v\n", r.Method, r.URL)
		log.Print("\tHeaders:\n")
		for key, val := range r.Header {
			log.Printf("\t\t%v: %v\n", key, val)
		}
		bodyBytes, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		if err == nil {
			log.Printf("\tBody: %v\n", string(bodyBytes))
		}
		next.ServeHTTP(w, r)
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/snorman7384/recipe-wizard/domain"
	"github.com/snorman7384/recipe-wizard/domerr"
)

type ContextKey string

const ContextUserKey ContextKey = "user-key"

type Config struct {
	Domain    domain.Config
	JwtSecret []byte
	Port      string
}

func (c *Config) Serve() {
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

	v1.Get("/grocery-lists/{grocery_list_id}/items", c.middlewareExtractUser(c.handleGetItemsForGroceryList()))
	v1.Post("/grocery-lists/{grocery_list_id}/items", c.middlewareExtractUser(c.handlePostItem()))
	v1.Get("/grocery-lists/{grocery_list_id}/items/{item_name}", c.middlewareExtractUser(c.handleGetItemsForGroceryListByName()))

	v1.Get("/items/{item_id}", c.middlewareExtractUser(c.handleGetItem()))
	v1.Put("/items/{item_id}/status", c.middlewareExtractUser(c.handleMarkItemStatus()))

	v1.Post("/users", c.handlePostUser())
	v1.Post("/login", c.handleLogin())

	server := &http.Server{
		Addr:              "0.0.0.0:" + c.Port,
		Handler:           r,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("Listening on port: ", c.Port)
	log.Fatal(server.ListenAndServe())
}

func (c *Config) handleIndex() http.HandlerFunc {
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

func (c *Config) handlePing() http.HandlerFunc {
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

func respondWithDomainError(w http.ResponseWriter, err error) {

	e, ok := err.(*domerr.DomainError)

	if !ok {
		log.Println("UNTYPED_DOMAIN_ERROR:", err.Error())
		e = domerr.ErrInternal
	}

	var code int
	switch e.Type() {
	case domerr.NotFound:
		code = http.StatusNotFound
	case domerr.UserNotFound:
		code = http.StatusUnauthorized
	case domerr.RecipeScraperFailure:
		code = http.StatusBadRequest
	case domerr.Forbidden:
		code = http.StatusForbidden
	case domerr.Internal:
		code = http.StatusInternalServerError
	default:
		code = http.StatusInternalServerError
	}

	respondWithError(w, code, e.Error())
}

func respondWithError(w http.ResponseWriter, code int, err string) {
	type response struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, response{Error: err})
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

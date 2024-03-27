package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/snorman7384/recipe-wizard/domain"
	"github.com/snorman7384/recipe-wizard/domerr"
	"golang.org/x/crypto/bcrypt"
)

type ContextKey string

const ContextUserKey ContextKey = "user-key"

type Config struct {
	Domain    domain.Config
	JwtSecret []byte
	Port      string
	Templates *template.Template
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

	r.Use(func(next http.Handler) http.Handler { return middlewareLogRequest(next) })

	r.Get("/login", c.handleLoginWWW())
	r.Post("/login", c.handlePostLoginWWW())

	r.Get("/", c.middlewareExtractUserFromCookie(c.handleIndexWWW()))
	r.Get("/recipes", c.middlewareExtractUserFromCookie(c.handleRecipesWWW()))
	r.Post("/recipes", c.middlewareExtractUserFromCookie(c.handlePostRecipesWWW()))
	r.Get("/recipes/{recipe_id}", c.middlewareExtractUserFromCookie(c.handleRecipeWWW()))

	r.Mount("/static", http.StripPrefix("/static", http.FileServer(http.Dir("/home/snorm/dev/recipe-wizard/static"))))

	v1 := chi.NewRouter()
	r.Mount("/v1", v1)

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

func (c *Config) handleIndexWWW() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.Templates.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			log.Println("Could not respond to index request: ", err)
			w.WriteHeader(500)
		}
	}
}

func (c *Config) handleRecipesWWW() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		recipes, err := c.Domain.GetRecipesForUser(r.Context(), user)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		err = c.Templates.ExecuteTemplate(w, "recipes.gohtml", recipes)
		if err != nil {
			log.Println("Could not respond to index request: ", err)
			w.WriteHeader(500)
		}
	}
}

func (c *Config) handleRecipeWWW() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		idString := chi.URLParam(r, "recipe_id")

		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Id is not an integer")
			return
		}

		recipe, err := c.Domain.GetRecipe(r.Context(), user, id)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		ingredients, err := c.Domain.GetIngredientsForRecipe(r.Context(), user, recipe)

		err = c.Templates.ExecuteTemplate(w, "recipe.gohtml", struct {
			domain.Recipe
			Ingredients []domain.Ingredient
		}{
			Recipe:      recipe,
			Ingredients: ingredients,
		})
		if err != nil {
			log.Println("Could not respond to index request: ", err)
			w.WriteHeader(500)
		}
	}
}

func (c *Config) handlePostRecipesWWW() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ContextUserKey).(domain.User)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve user")
			return
		}

		url := r.FormValue("url")

		_, err := c.Domain.CreateRecipeFromUrl(r.Context(), user, url)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		http.Redirect(w, r, "/recipes", 302)
	}
}

func (c *Config) handleLoginWWW() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		previousError := r.URL.Query().Has("invalid-form")
		err := c.Templates.ExecuteTemplate(w, "login.gohtml", struct{ PreviousError bool }{PreviousError: previousError})

		if err != nil {
			log.Println("Could not respond to index request: ", err)
			w.WriteHeader(500)
		}
	}
}

func (c *Config) handlePostLoginWWW() http.HandlerFunc {
	jwtDuration := time.Hour * 24
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			http.Redirect(w, r, "/login?invalid-form", 302)
			return
		}

		user, err := c.Domain.GetUserByUsername(r.Context(), username)
		if err != nil {
			respondWithDomainError(w, err) // TODO:
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Password") // TODO:
			return
		}

		issuedAt := jwt.NewNumericDate(time.Now())
		expiresAt := jwt.NewNumericDate(issuedAt.Add(jwtDuration))

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "recipe-wizard",
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAt,
			Subject:   fmt.Sprint(user.ID),
		})

		tokenString, err := token.SignedString(c.JwtSecret)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error()) // TODO:
			return
		}

		w.Header().Add("Set-Cookie", fmt.Sprintf("userAccessToken=%s; HttpOnly", tokenString))

		http.Redirect(w, r, "/", 302)
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

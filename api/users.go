package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snorman7384/recipe-wizard/domain"
	"golang.org/x/crypto/bcrypt"
)

const maxJwtDuration = time.Hour * 24

func (c *Config) handlePostUser() http.HandlerFunc {

	type request struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		FirstName        string `json:"first_name"`
		LastName         string `json:"last_name"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		FirstName string    `json:"first_name,omitempty"`
		LastName  string    `json:"last_name,omitempty"`
		Token     string    `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			respondWithError(w, http.StatusBadRequest, "Password is too long")
			return
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := c.Domain.CreateUser(r.Context(), domain.CreateUserParams{
			Username:       req.Username,
			HashedPassword: string(hashedPasswordBytes),
			FirstName:      req.FirstName,
			LastName:       req.LastName,
		})
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		jwtDuration := time.Second * time.Duration(req.ExpiresInSeconds)
		if req.ExpiresInSeconds <= 0 || jwtDuration > maxJwtDuration {
			jwtDuration = maxJwtDuration
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
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		res := response{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Token:     tokenString,
		}

		respondWithJSON(w, http.StatusCreated, &res)
	}
}

func (c *Config) handleLogin() http.HandlerFunc {
	type request struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		Token string `json:"token"`
	}

	maxJwtDuration := time.Hour * 24

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := c.Domain.GetUserByUsername(r.Context(), req.Username)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect Password")
			return
		}

		jwtDuration := time.Second * time.Duration(req.ExpiresInSeconds)
		if req.ExpiresInSeconds <= 0 || jwtDuration > maxJwtDuration {
			jwtDuration = maxJwtDuration
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
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		res := response{
			Token: tokenString,
		}

		respondWithJSON(w, http.StatusOK, &res)
	}
}

func (c *Config) middlewareExtractUser(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authString := r.Header.Get("Authorization")
		if authString == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing Authorization Header")
			return
		}

		re := regexp.MustCompile(`\s+`)
		authList := re.Split(authString, -1)

		if len(authList) != 2 || authList[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Malformed Authorization Header")
			return
		}

		tokenString := authList[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(c.JwtSecret), nil
		})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userIdString, err := token.Claims.GetSubject()
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "No user specified")
			return
		}

		userId, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid user id")
			return
		}

		user, err := c.Domain.GetUser(r.Context(), userId)
		if err != nil {
			respondWithDomainError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func (c *Config) middlewareExtractUserFromCookie(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("userAccessToken")
		if err != nil {
			fmt.Println("bad cookie:", err.Error())
			http.Redirect(w, r, "/login", 302)
			return
		}
		if cookie.Valid() != nil {
			fmt.Println("invalid cookie")
			http.Redirect(w, r, "/login", 302)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			return []byte(c.JwtSecret), nil
		})
		if err != nil {
			fmt.Println("couln't parse JWT")
			http.Redirect(w, r, "/login", 302)
			return
		}

		userIdString, err := token.Claims.GetSubject()
		if err != nil {
			fmt.Println("no subject")
			http.Redirect(w, r, "/login", 302)
			return
		}

		userId, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			fmt.Println("id not number")
			http.Redirect(w, r, "/login", 302)
			return
		}

		user, err := c.Domain.GetUser(r.Context(), userId)
		if err != nil {
			respondWithDomainError(w, err) // TODO:
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

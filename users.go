package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/snorman7384/recipe-wizard/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (c *config) handlePostUser() http.HandlerFunc {

	maxJwtDuration := time.Hour * 24

	type request struct {
		Username         string  `json:"username"`
		Password         string  `json:"password"`
		FirstName        *string `json:"first_name"`
		LastName         *string `json:"last_name"`
		ExpiresInSeconds int     `json:"expires_in_seconds"`
	}

	type response struct {
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		FirstName *string   `json:"first_name,omitempty"`
		LastName  *string   `json:"last_name,omitempty"`
		Token     string    `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		firstName := sqlNullStringFromStringPointer(req.FirstName)
		lastName := sqlNullStringFromStringPointer(req.LastName)

		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			respondWithError(w, http.StatusBadRequest, "Password is too long")
			return
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		now := time.Now()

		err = c.DB.CreateUser(r.Context(), database.CreateUserParams{
			CreatedAt:      now,
			UpdatedAt:      now,
			Username:       req.Username,
			HashedPassword: string(hashedPasswordBytes),
			FirstName:      firstName,
			LastName:       lastName,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		id, err := c.DB.GetLastInsertID(r.Context())

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		dbUser, err := c.DB.GetUser(r.Context(), id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
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
			Subject:   fmt.Sprint(id),
		})

		tokenString, err := token.SignedString(c.JwtSecret)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		res := response{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Username:  dbUser.Username,
			FirstName: stringPointerFromSqlNullString(dbUser.FirstName),
			LastName:  stringPointerFromSqlNullString(dbUser.LastName),
			Token:     tokenString,
		}

		respondWithJSON(w, http.StatusCreated, &res)
	}
}

func (c *config) handleLogin() http.HandlerFunc {
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

		user, err := c.DB.GetUserByUsername(r.Context(), req.Username)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
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

package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload loginRequest
	if err := readJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := app.store.Users.GetByUsername(r.Context(), payload.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	// Same 401 for "not found" and "wrong password" to prevent user enumeration.
	if user == nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     now.Add(72 * time.Hour).Unix(),
		"iat":     now.Unix(),
		"iss":     "blog-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(app.config.jwtSecret))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenStr}) //nolint:errcheck
}

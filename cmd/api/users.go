package main

import (
	"net/http"

	"github.com/gustavooarantes/Blog/internal/store"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserRequest
	if err := readJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		writeError(w, http.StatusUnprocessableEntity, "username, email, and password are required")
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	writeJSON(w, http.StatusCreated, user) //nolint:errcheck
}

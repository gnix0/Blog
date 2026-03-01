package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gustavooarantes/Blog/internal/store"
)

func getPostIDParam(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title         string   `json:"title"`
		Content       string   `json:"content"`
		Tags          []string `json:"tags"`
		CoverImageURL string   `json:"cover_image_url"`
	}
	if err := readJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := r.Context().Value(userIDKey).(int64)

	post := &store.Post{
		Title:         payload.Title,
		Content:       payload.Content,
		Tags:          payload.Tags,
		CoverImageURL: payload.CoverImageURL,
		UserID:        userID,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create post")
		return
	}

	writeJSON(w, http.StatusCreated, post) //nolint:errcheck
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getPostIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if post == nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}

	writeJSON(w, http.StatusOK, post) //nolint:errcheck
}

func (app *application) listPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.store.Posts.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if posts == nil {
		posts = []store.Post{}
	}

	writeJSON(w, http.StatusOK, posts) //nolint:errcheck
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getPostIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	existing, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}

	var payload struct {
		Title         *string  `json:"title"`
		Content       *string  `json:"content"`
		Tags          []string `json:"tags"`
		CoverImageURL *string  `json:"cover_image_url"`
	}
	if err := readJSON(w, r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if payload.Title != nil {
		existing.Title = *payload.Title
	}
	if payload.Content != nil {
		existing.Content = *payload.Content
	}
	if payload.Tags != nil {
		existing.Tags = payload.Tags
	}
	if payload.CoverImageURL != nil {
		existing.CoverImageURL = *payload.CoverImageURL
	}

	if err := app.store.Posts.Update(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, "could not update post")
		return
	}

	writeJSON(w, http.StatusOK, existing) //nolint:errcheck
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getPostIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid post ID")
		return
	}

	err = app.store.Posts.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete post")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gustavooarantes/Blog/internal/store"
)

type application struct {
	config config
	logger *log.Logger
	store  store.Storage
}

type cloudinaryConfig struct {
	cloudName string
	apiKey    string
	apiSecret string
}

type config struct {
	addr       string
	jwtSecret  string
	uploadsDir string
	cloudinary cloudinaryConfig
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Serve locally-uploaded files (dev only; in production Cloudinary returns absolute URLs)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(app.config.uploadsDir))))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		// Public
		r.Post("/users", app.registerUserHandler)
		r.Post("/auth/token", app.createTokenHandler)
		r.Get("/posts", app.listPostsHandler)
		r.Get("/posts/{postID}", app.getPostHandler)

		// Protected
		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware)

			r.Post("/posts", app.createPostHandler)
			r.Patch("/posts/{postID}", app.updatePostHandler)
			r.Delete("/posts/{postID}", app.deletePostHandler)
			r.Post("/uploads", app.uploadImageHandler)
		})
	})

	// Serve the React SPA — must be registered last.
	// Any path that isn't a real file in web/dist/ falls back to index.html
	// so client-side routing (React Router) works correctly.
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "web/dist" + r.URL.Path
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			http.ServeFile(w, r, "web/dist/index.html")
			return
		}
		http.FileServer(http.Dir("web/dist")).ServeHTTP(w, r)
	}))

	return r
}

func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server started at %s", app.config.addr)

	return srv.ListenAndServe()
}

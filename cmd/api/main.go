package main

import (
	"log"
	"os"

	"github.com/gustavooarantes/Blog/internal/db"
	"github.com/gustavooarantes/Blog/internal/env"
	"github.com/gustavooarantes/Blog/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for local development; in production env vars are injected by the runtime.
	_ = godotenv.Load()

	// Heroku sets PORT dynamically; local dev falls back to ADDR.
	addr := env.GetString("ADDR", ":8080")
	if port := env.GetString("PORT", ""); port != "" {
		addr = ":" + port
	}

	// Heroku Postgres sets DATABASE_URL automatically; local dev uses DB_ADDR.
	dbAddr := env.GetString("DATABASE_URL", env.GetString("DB_ADDR", ""))
	if dbAddr == "" {
		log.Fatal("DB_ADDR (or DATABASE_URL) environment variable is required")
	}

	cfg := config{
		addr:       addr,
		jwtSecret:  env.GetString("JWT_SECRET", ""),
		uploadsDir: env.GetString("UPLOADS_DIR", "/uploads"),
		cloudinary: cloudinaryConfig{
			cloudName: env.GetString("CLOUDINARY_CLOUD_NAME", ""),
			apiKey:    env.GetString("CLOUDINARY_API_KEY", ""),
			apiSecret: env.GetString("CLOUDINARY_API_SECRET", ""),
		},
	}

	if cfg.jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	database, err := db.Open(dbAddr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	storage := store.NewPostgresStorage(database)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config: cfg,
		logger: logger,
		store:  storage,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

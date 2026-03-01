package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/gustavooarantes/Blog/internal/env"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env for local development; in production env vars are injected by the runtime.
	_ = godotenv.Load()

	// Heroku Postgres sets DATABASE_URL automatically; local dev uses DB_ADDR.
	dbAddr := env.GetString("DATABASE_URL", env.GetString("DB_ADDR", ""))
	if dbAddr == "" {
		log.Fatal("DB_ADDR (or DATABASE_URL) environment variable is required")
	}

	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	matches, err := filepath.Glob("cmd/migrate/migrations/*.up.sql")
	if err != nil {
		log.Fatalf("glob migrations: %v", err)
	}
	sort.Strings(matches)

	for _, path := range matches {
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read %s: %v", path, err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("exec %s: %v", path, err)
		}
		log.Printf("applied: %s", path)
	}

	log.Println("migrations complete")
}

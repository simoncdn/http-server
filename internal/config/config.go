package config

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/simoncdn/http-server/internal/database"
)

type Config struct {
	Port           string
	StaticDir      string
	FileserverHits atomic.Int32
	DB             *database.Queries
}

func New() *Config {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to databse:", err)
	}
	dbQueries := database.New(db)

	return &Config{
		Port:           "8080",
		StaticDir:      "./web/static",
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
	}
}

func (c *Config) GetHits() int32 {
	return c.FileserverHits.Load()
}

func (c *Config) IncrementHits() {
	c.FileserverHits.Add(1)
}

func (c *Config) ResetHits() {
	c.FileserverHits.Store(0)
}

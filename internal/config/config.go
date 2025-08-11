package config

import (
	"sync/atomic"
)

type Config struct {
	Port           string
	StaticDir      string
	FileserverHits atomic.Int32
}

func New() *Config {
	return &Config{
		Port:           "8080",
		StaticDir:      "./web/static",
		FileserverHits: atomic.Int32{},
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

package config

import (
	"os"
	"path/filepath"
)

type Config struct{ DataDir, DatabasePath, ListenAddr string }

func Load() Config {
	dir := os.Getenv("HYPERION_DATA_DIR")
	if dir == "" {
		dir = "./data"
	}
	addr := os.Getenv("HYPERION_LISTEN_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	return Config{DataDir: dir, DatabasePath: filepath.Join(dir, "hyperion.db"), ListenAddr: addr}
}

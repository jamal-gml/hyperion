package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jamal-gml/hyperion/internal/config"
	"github.com/jamal-gml/hyperion/internal/httpapi"
	"github.com/jamal-gml/hyperion/internal/mcpserver"
	"github.com/jamal-gml/hyperion/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cfg := config.Load()
	switch os.Args[1] {
	case "serve":
		if err := serve(cfg); err != nil {
			slog.Error("server stopped", "error", err)
			os.Exit(1)
		}
	case "mcp":
		if len(os.Args) < 3 || os.Args[2] != "stdio" {
			usage()
			os.Exit(2)
		}
		if err := runStdio(cfg); err != nil {
			slog.Error("stdio MCP stopped", "error", err)
			os.Exit(1)
		}
	case "version":
		fmt.Println("Hyperion dev")
	default:
		usage()
		os.Exit(2)
	}
}

func serve(cfg config.Config) error {
	db, err := store.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer db.Close()
	api := httpapi.New(db, cfg)
	return api.Run()
}

func runStdio(cfg config.Config) error {
	db, err := store.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer db.Close()
	id := mcpserver.IdentityFromEnv()
	server := mcpserver.New(db, id)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return server.Run(ctx)
}

func usage() { fmt.Fprintln(os.Stderr, "Usage: hyperion <serve|mcp stdio|version>") }

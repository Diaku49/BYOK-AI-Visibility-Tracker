package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Diaku49/AI-visibility-tracker/config"
	"github.com/Diaku49/AI-visibility-tracker/internal/api"
	"github.com/Diaku49/AI-visibility-tracker/internal/api/handlers"
	"github.com/Diaku49/AI-visibility-tracker/internal/logger"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
	"github.com/Diaku49/AI-visibility-tracker/internal/worker"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("load config", "error", err)
		os.Exit(1)
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	keyCipher, err := pkg.NewKeyCipherFromBase64(cfg.MasterKey)
	if err != nil {
		log.Error("create key cipher", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	st, err := store.Connect(ctx, databaseURL)
	if err != nil {
		log.Error("connect to database", "error", err)
		os.Exit(1)
	}

	server := api.NewServer(st, cfg)
	handler := handlers.NewServerHandler(st, cfg, keyCipher)
	coordinator := worker.NewCoordinator(st, log, keyCipher)

	go coordinator.Start()

	address := ":" + cfg.Port
	log.Info("API server started", "address", address)
	if err := http.ListenAndServe(address, server.Route(handler)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("API server stopped", "error", err)
	}
}

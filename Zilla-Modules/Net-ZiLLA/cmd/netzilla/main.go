package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"net-zilla/internal/api"
	"net-zilla/internal/config"
	"net-zilla/internal/services"
	"net-zilla/internal/storage"
	"net-zilla/internal/utils"
	"net-zilla/pkg/logger"
)

func main() {
	// 1. Config & Logger
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	l := logger.NewLogger()

	// 2. Storage
	db, err := storage.NewDatabase("netzilla.db")
	if err != nil {
		l.Error("Failed to initialize database: %v", err)
	} else {
		defer db.Close()
	}

	// 3. Core Services (New Architecture)
	analysisService := services.NewAnalysisService(l, db, cfg)

	// 4. Lifecycle Management
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		l.Info("🛑 Shutting down...")
		cancel()
	}()

	// 5. Entry Point Selection
	if cfg.Server.EnableAPI {
		l.Info("Starting Net-Zilla API server...")
		apiServer := api.NewServer(analysisService, l, cfg)
		if err := apiServer.Run(ctx); err != nil {
			l.Error("API server failed: %v", err)
			os.Exit(1)
		}
	} else {
		l.Info("Starting Net-Zilla CLI...")
		menu := utils.NewMenu(analysisService, l)
		if err := menu.Run(); err != nil {
			l.Error("CLI menu failed: %v", err)
			os.Exit(1)
		}
	}
}

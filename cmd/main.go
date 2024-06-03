package main

import (
	"IPFS-CLUSTER-MANAGER/internal/api"
	"IPFS-CLUSTER-MANAGER/internal/core/config"
	"IPFS-CLUSTER-MANAGER/internal/core/services"
	"IPFS-CLUSTER-MANAGER/internal/log"
	"context"
	"fmt"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error(context.Background(), "Error loading configuration: %s", err)
		return
	}

	ctx, cancel := context.WithCancel(log.NewContext(context.Background(), cfg.Log.Level, cfg.Log.Mode, os.Stdout))
	defer cancel()

	ipfsService := services.NewIpfs(cfg)

	server := api.NewServer(ipfsService)

	strictHandler := api.NewStrictHandler(server, nil)
	handler := api.Handler(strictHandler)

	log.Info(ctx, fmt.Sprintf("Starting server on %s:%d", cfg.ServerUrl, cfg.ServerPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), handler); err != nil {
		log.Error(ctx, "Error starting server: %s", err)
	}
}

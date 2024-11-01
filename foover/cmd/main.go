package main

import (
	"context"
	"fmt"
	golog "log"
	"net/http"
	"time"

	_ "foover/docs"
	"foover/internal/config"
	"foover/internal/log"
	"foover/internal/service"
	"foover/internal/store/mongo"
	httpTransport "foover/internal/transport/http"
)

// @title Foover API
// @version 1.0
// @description This is the API documentation for Foover.

// @host localhost:8080
// @BasePath /
func main() {

	// Load configuration
	cfg, err := config.LoadEnvVars()
	if err != nil {
		golog.Fatalf("Failed to load configuration: %v", err)
	}

	// initialize logger
	logger := log.InitializeLogger(cfg.Service.LogLevel)

	// Initialize MongoDB store
	store, err := mongo.NewStore(cfg.Mongo)
	if err != nil {
		logger.Error("Failed to initialize MongoDB store: %v", err)
	}
	defer store.Close()

	// Initialize services
	sessionService := service.NewSessionService(store)
	voteService := service.NewVoteService(store)
	aggregationService := service.NewAggregationService(store)
	productService := service.NewProductService(store)

	// Initialize HTTP server
	router := httpTransport.NewRouter(sessionService, voteService, aggregationService, productService, logger)

	// Fetch and store products
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = productService.FetchAndStoreProducts(ctx, cfg.ExternalAPI.ProductAPIURL)
	if err != nil {
		logger.Error("Failed to fetch and store products: %v", err)
	}
	logger.Info("Products fetched and stored successfully")

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	logger.Info(fmt.Sprintf("Server is running on %s", cfg.HTTPServer.Address))
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}

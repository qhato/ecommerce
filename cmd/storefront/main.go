package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/config"

	// Catalog
	catalogQueries "github.com/qhato/ecommerce/internal/catalog/application/queries"
	catalogPersistence "github.com/qhato/ecommerce/internal/catalog/infrastructure/persistence"
	catalogHttp "github.com/qhato/ecommerce/internal/catalog/ports/http"

	// Customer
	customerCommands "github.com/qhato/ecommerce/internal/customer/application/commands"
	customerQueries "github.com/qhato/ecommerce/internal/customer/application/queries"
	customerPersistence "github.com/qhato/ecommerce/internal/customer/infrastructure/persistence"
	customerHttp "github.com/qhato/ecommerce/internal/customer/ports/http"

	// Order
	orderQueries "github.com/qhato/ecommerce/internal/order/application/queries"
	orderPersistence "github.com/qhato/ecommerce/internal/order/infrastructure/persistence"
	orderHttp "github.com/qhato/ecommerce/internal/order/ports/http"

	// Fulfillment
	fulfillmentPersistence "github.com/qhato/ecommerce/internal/fulfillment/infrastructure/persistence"
	fulfillmentHttp "github.com/qhato/ecommerce/internal/fulfillment/ports/http"

	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/middleware"
	"github.com/qhato/ecommerce/pkg/validator"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	log.Info("Starting Storefront API server", "version", "1.0.0")

	// Initialize database (read-mostly connection pool for storefront)
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()
	log.Info("Connected to database")

	// Initialize cache (important for storefront performance)
	var cacheStore cache.Cache
	if cfg.Cache.Type == "redis" {
		cacheStore, err = cache.NewRedisCache(cfg.Cache.Redis)
		if err != nil {
			log.Fatal("Failed to connect to Redis", "error", err)
		}
		log.Info("Connected to Redis cache")
	} else {
		cacheStore = cache.NewMemoryCache()
		log.Info("Using in-memory cache")
	}

	// Initialize event bus (for customer registration, etc.)
	eventBus := event.NewMemoryBus(log)
	log.Info("Event bus initialized")

	// Initialize validator
	val := validator.New()

	// ========== CATALOG BOUNDED CONTEXT ==========

	// Catalog repositories
	productRepo := catalogPersistence.NewPostgresProductRepository(db)
	categoryRepo := catalogPersistence.NewPostgresCategoryRepository(db)
	skuRepo := catalogPersistence.NewPostgresSKURepository(db)

	// Catalog query handlers (storefront is mostly read-only)
	productQueryHandler := catalogQueries.NewProductQueryHandler(productRepo, cacheStore, log)
	categoryQueryHandler := catalogQueries.NewCategoryQueryHandler(categoryRepo, cacheStore, log)
	skuQueryHandler := catalogQueries.NewSKUQueryHandler(skuRepo, cacheStore, log)

	// Catalog HTTP handlers
	storefrontCatalogHandler := catalogHttp.NewStorefrontCatalogHandler(productQueryHandler, categoryQueryHandler, skuQueryHandler, log)

	// ========== CUSTOMER BOUNDED CONTEXT ==========

	// Customer repositories
	customerRepo := customerPersistence.NewPostgresCustomerRepository(db)

	// Customer command handlers (for registration)
	customerCommandHandler := customerCommands.NewCustomerCommandHandler(customerRepo, eventBus, log)

	// Customer query handlers
	customerQueryHandler := customerQueries.NewCustomerQueryHandler(customerRepo, cacheStore, log)

	// Customer HTTP handlers
	storefrontCustomerHandler := customerHttp.NewStorefrontCustomerHandler(customerCommandHandler, customerQueryHandler, val, log)

	// ========== ORDER BOUNDED CONTEXT ==========

	// Order repositories
	orderRepo := orderPersistence.NewPostgresOrderRepository(db)

	// Order query handlers
	orderQueryHandler := orderQueries.NewOrderQueryHandler(orderRepo, cacheStore, log)

	// Order HTTP handlers
	storefrontOrderHandler := orderHttp.NewStorefrontOrderHandler(orderQueryHandler, log)

	// ========== FULFILLMENT BOUNDED CONTEXT ==========

	// Fulfillment repositories
	shipmentRepo := fulfillmentPersistence.NewPostgresShipmentRepository(db)

	// Fulfillment HTTP handlers
	storefrontShipmentHandler := fulfillmentHttp.NewStorefrontShipmentHandler(shipmentRepo, log)

	// ========== ROUTER SETUP ==========

	// Setup router
	r := chi.NewRouter()

	// Apply global middleware
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.Recoverer(log))
	r.Use(middleware.CORS())

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// API info
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"service": "E-Commerce Storefront API",
			"version": "1.0.0",
			"description": "Public API for e-commerce storefront including catalog, orders, and customer management"
		}`))
	})

	// Register storefront routes (public, some may require auth in production)
	storefrontCatalogHandler.RegisterRoutes(r)
	storefrontCustomerHandler.RegisterRoutes(r)
	storefrontOrderHandler.RegisterRoutes(r)
	storefrontShipmentHandler.RegisterRoutes(r)

	log.Info("All storefront contexts initialized",
		"contexts", "catalog, customer, order, fulfillment")

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Storefront API server listening", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Storefront API server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Storefront API server stopped")
}

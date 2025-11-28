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
	catalogCommands "github.com/qhato/ecommerce/internal/catalog/application/commands"
	catalogQueries "github.com/qhato/ecommerce/internal/catalog/application/queries"
	catalogPersistence "github.com/qhato/ecommerce/internal/catalog/infrastructure/persistence"
	catalogHttp "github.com/qhato/ecommerce/internal/catalog/ports/http"

	// Customer
	customerCommands "github.com/qhato/ecommerce/internal/customer/application/commands"
	customerQueries "github.com/qhato/ecommerce/internal/customer/application/queries"
	customerPersistence "github.com/qhato/ecommerce/internal/customer/infrastructure/persistence"
	customerHttp "github.com/qhato/ecommerce/internal/customer/ports/http"

	// Order
	orderCommands "github.com/qhato/ecommerce/internal/order/application/commands"
	orderQueries "github.com/qhato/ecommerce/internal/order/application/queries"
	orderPersistence "github.com/qhato/ecommerce/internal/order/infrastructure/persistence"
	orderHttp "github.com/qhato/ecommerce/internal/order/ports/http"

	// Payment
	paymentCommands "github.com/qhato/ecommerce/internal/payment/application/commands"
	paymentQueries "github.com/qhato/ecommerce/internal/payment/application/queries"
	paymentPersistence "github.com/qhato/ecommerce/internal/payment/infrastructure/persistence"
	paymentHttp "github.com/qhato/ecommerce/internal/payment/ports/http"

	// Fulfillment
	fulfillmentCommands "github.com/qhato/ecommerce/internal/fulfillment/application/commands"
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
	log.Info("Starting Admin API server", "version", "1.0.0")

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()
	log.Info("Connected to database")

	// Initialize cache
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

	// Initialize event bus
	eventBus := event.NewMemoryBus(log)
	log.Info("Event bus initialized")

	// Initialize validator
	val := validator.New()

	// ========== CATALOG BOUNDED CONTEXT ==========

	// Catalog repositories
	productRepo := catalogPersistence.NewPostgresProductRepository(db)
	categoryRepo := catalogPersistence.NewPostgresCategoryRepository(db)
	skuRepo := catalogPersistence.NewPostgresSKURepository(db)

	// Catalog command handlers
	productCommandHandler := catalogCommands.NewProductCommandHandler(productRepo, eventBus, val, log)
	categoryCommandHandler := catalogCommands.NewCategoryCommandHandler(categoryRepo, eventBus, val, log)
	skuCommandHandler := catalogCommands.NewSKUCommandHandler(skuRepo, eventBus, val, log)

	// Catalog query handlers
	productQueryHandler := catalogQueries.NewProductQueryHandler(productRepo, cacheStore, log)
	categoryQueryHandler := catalogQueries.NewCategoryQueryHandler(categoryRepo, cacheStore, log)
	skuQueryHandler := catalogQueries.NewSKUQueryHandler(skuRepo, cacheStore, log)

	// Catalog HTTP handlers
	adminProductHandler := catalogHttp.NewAdminProductHandler(productCommandHandler, productQueryHandler, log)
	adminCategoryHandler := catalogHttp.NewAdminCategoryHandler(categoryCommandHandler, categoryQueryHandler, log)
	adminSKUHandler := catalogHttp.NewAdminSKUHandler(skuCommandHandler, skuQueryHandler, log)

	// ========== CUSTOMER BOUNDED CONTEXT ==========

	// Customer repositories
	customerRepo := customerPersistence.NewPostgresCustomerRepository(db)

	// Customer command handlers
	customerCommandHandler := customerCommands.NewCustomerCommandHandler(customerRepo, eventBus, log)

	// Customer query handlers
	customerQueryHandler := customerQueries.NewCustomerQueryHandler(customerRepo, cacheStore, log)

	// Customer HTTP handlers
	adminCustomerHandler := customerHttp.NewAdminCustomerHandler(customerCommandHandler, customerQueryHandler, val, log)

	// ========== ORDER BOUNDED CONTEXT ==========

	// Order repositories
	orderRepo := orderPersistence.NewPostgresOrderRepository(db)

	// Order command handlers
	orderCommandHandler := orderCommands.NewOrderCommandHandler(orderRepo, eventBus, log)

	// Order query handlers
	orderQueryHandler := orderQueries.NewOrderQueryHandler(orderRepo, cacheStore, log)

	// Order HTTP handlers
	adminOrderHandler := orderHttp.NewAdminOrderHandler(orderCommandHandler, orderQueryHandler, val, log)

	// ========== PAYMENT BOUNDED CONTEXT ==========

	// Payment repositories
	paymentRepo := paymentPersistence.NewPostgresPaymentRepository(db)

	// Payment command handlers
	paymentCommandHandler := paymentCommands.NewPaymentCommandHandler(paymentRepo, eventBus, log)

	// Payment query handlers
	paymentQueryHandler := paymentQueries.NewPaymentQueryHandler(paymentRepo, cacheStore, log)

	// Payment HTTP handlers
	adminPaymentHandler := paymentHttp.NewAdminPaymentHandler(paymentCommandHandler, paymentQueryHandler, val, log)

	// ========== FULFILLMENT BOUNDED CONTEXT ==========

	// Fulfillment repositories
	shipmentRepo := fulfillmentPersistence.NewPostgresShipmentRepository(db)

	// Fulfillment command handlers
	shipmentCommandHandler := fulfillmentCommands.NewShipmentCommandHandler(shipmentRepo, eventBus, log)

	// Fulfillment HTTP handlers
	adminShipmentHandler := fulfillmentHttp.NewAdminShipmentHandler(shipmentCommandHandler, shipmentRepo, val, log)

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

	// Register routes (protected with auth middleware for production)
	// For now, routes are open. In production, add: r.Use(middleware.Auth(jwtSecret))

	// Catalog routes
	adminProductHandler.RegisterRoutes(r)
	adminCategoryHandler.RegisterRoutes(r)
	adminSKUHandler.RegisterRoutes(r)

	// Customer routes
	adminCustomerHandler.RegisterRoutes(r)

	// Order routes
	adminOrderHandler.RegisterRoutes(r)

	// Payment routes
	adminPaymentHandler.RegisterRoutes(r)

	// Fulfillment routes
	adminShipmentHandler.RegisterRoutes(r)

	log.Info("All bounded contexts initialized",
		"contexts", "catalog, customer, order, payment, fulfillment")

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
		log.Info("Admin API server listening", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Admin API server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Admin API server stopped")
}

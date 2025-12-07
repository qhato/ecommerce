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
	catalogApp "github.com/qhato/ecommerce/internal/catalog/application"
	//catalogCommands "github.com/qhato/ecommerce/internal/catalog/application/commands"
	catalogQueries "github.com/qhato/ecommerce/internal/catalog/application/queries"
	catalogPersistence "github.com/qhato/ecommerce/internal/catalog/infrastructure/persistence"
	catalogHttp "github.com/qhato/ecommerce/internal/catalog/ports/http"

	// Customer
	customerCommands "github.com/qhato/ecommerce/internal/customer/application/commands"
	customerQueries "github.com/qhato/ecommerce/internal/customer/application/queries"
	customerPersistence "github.com/qhato/ecommerce/internal/customer/infrastructure/persistence"
	customerHttp "github.com/qhato/ecommerce/internal/customer/ports/http"

	// Order
	orderApp "github.com/qhato/ecommerce/internal/order/application"
	//orderCommands "github.com/qhato/ecommerce/internal/order/application/commands"
	orderQueries "github.com/qhato/ecommerce/internal/order/application/queries"
	orderPersistence "github.com/qhato/ecommerce/internal/order/infrastructure/persistence"
	orderHttp "github.com/qhato/ecommerce/internal/order/ports/http"

	// Offer
	offerApp "github.com/qhato/ecommerce/internal/offer/application"
	offerPersistence "github.com/qhato/ecommerce/internal/offer/infrastructure/persistence"

	// Inventory
	inventoryApp "github.com/qhato/ecommerce/internal/inventory/application"
	inventoryPersistence "github.com/qhato/ecommerce/internal/inventory/infrastructure/persistence"

	// Tax
	// taxApp "github.com/qhato/ecommerce/internal/tax/application" // Commented out - old tax implementation
	// taxPersistence "github.com/qhato/ecommerce/internal/tax/infrastructure/persistence" // Commented out - old tax implementation

	// Payment
	//paymentCommands "github.com/qhato/ecommerce/internal/payment/application/commands"
	//paymentQueries "github.com/qhato/ecommerce/internal/payment/application/queries"
	//paymentPersistence "github.com/qhato/ecommerce/internal/payment/infrastructure/persistence"
	//paymentHttp "github.com/qhato/ecommerce/internal/payment/ports/http"

	// Fulfillment
	//fulfillmentCommands "github.com/qhato/ecommerce/internal/fulfillment/application/commands"
	fulfillmentPersistence "github.com/qhato/ecommerce/internal/fulfillment/infrastructure/persistence"
	fulfillmentHttp "github.com/qhato/ecommerce/internal/fulfillment/ports/http"

	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/database"
	"github.com/qhato/ecommerce/pkg/event"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/middleware"
	"github.com/qhato/ecommerce/pkg/renderer"
	"github.com/qhato/ecommerce/pkg/validator"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	err = logger.Initialize(cfg.App.Environment, cfg.App.LogLevel)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log := logger.Get()
	log.WithField("version", cfg.App.Version).Info("Starting Storefront API server")

	// Initialize database (read-mostly connection pool for storefront)
	db, err := database.New(context.Background(), database.Config{ // Convert config.DatabaseConfig to database.Config
		Host: cfg.Database.Host,
		Port: cfg.Database.Port,
		User: cfg.Database.User,
		Password: cfg.Database.Password,
		Database: cfg.Database.Database,
		SSLMode: cfg.Database.SSLMode,
		MaxConnections: cfg.Database.MaxConnections,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxLifetime: cfg.Database.MaxLifetime,
		MaxIdleTime: cfg.Database.MaxIdleTime,
	})
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()
	log.Info("Connected to database")

	// Initialize cache (important for storefront performance)
	var cacheStore cache.Cache
	if cfg.Redis.Host != "" { // Check Redis host for cache type
		cacheStore, err = cache.NewRedisCache(cache.RedisConfig{ // Convert config.RedisConfig to cache.RedisConfig
			Host: cfg.Redis.Host,
			Port: cfg.Redis.Port,
			Password: cfg.Redis.Password,
			Database: cfg.Redis.Database,
			PoolSize: cfg.Redis.PoolSize,
			Prefix: "admin_api", // Assuming a prefix for admin cache
		})
		if err != nil {
			log.WithError(err).Fatal("Failed to connect to Redis")
		}
		log.Info("Connected to Redis cache")
	} else {
		cacheStore = cache.NewMemoryCache(cfg.Redis.TTL, cfg.Redis.TTL/2) // Provide arguments
		log.Info("Using in-memory cache")
	}

	// Initialize event bus (for customer registration, etc.)
	eventBus := event.NewMemoryBus()
	log.Info("Event bus initialized")

	// Initialize validator
	val := validator.New()

	// Initialize template renderer
	tmplRenderer, err := renderer.NewTemplateRenderer("web", log)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize template renderer")
	}
	log.Info("Template renderer initialized")

	// ========== CATALOG BOUNDED CONTEXT ==========

	// Catalog repositories
	productRepo := catalogPersistence.NewPostgresProductRepository(db)
	productAttributeRepo := catalogPersistence.NewPostgresProductAttributeRepository(db)
	categoryRepo := catalogPersistence.NewPostgresCategoryRepository(db)
	categoryAttributeRepo := catalogPersistence.NewPostgresCategoryAttributeRepository(db)
	skuRepo := catalogPersistence.NewPostgresSKURepository(db)
	productOptionXrefRepo := catalogPersistence.NewPostgresProductOptionXrefRepository(db)
	categoryProductXrefRepo := catalogPersistence.NewPostgresCategoryProductXrefRepository(db)
	skuAttributeRepo := catalogPersistence.NewPostgresSKUAttributeRepository(db)
	skuProductOptionValueXrefRepo := catalogPersistence.NewPostgresSkuProductOptionValueXrefRepository(db)
	productOptionRepo := catalogPersistence.NewPostgresProductOptionRepository(db)
	productOptionValueRepo := catalogPersistence.NewPostgresProductOptionValueRepository(db)

	// Catalog application services
	productService := catalogApp.NewProductService(productRepo, productAttributeRepo, productOptionXrefRepo, categoryProductXrefRepo)
	_ = catalogApp.NewCategoryService(categoryRepo, categoryAttributeRepo) // Assigned to _
	skuService := catalogApp.NewSkuService(skuRepo, skuAttributeRepo, skuProductOptionValueXrefRepo)
	_ = catalogApp.NewProductOptionService(productOptionRepo, productOptionValueRepo) // Assigned to _

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
	customerCommandHandler := customerCommands.NewCustomerCommandHandler(customerRepo, eventBus, val, log)

	// Customer query handlers
	customerQueryHandler := customerQueries.NewCustomerQueryHandler(customerRepo, cacheStore, log)

	// Customer HTTP handlers
	storefrontCustomerHandler := customerHttp.NewStorefrontCustomerHandler(customerCommandHandler, customerQueryHandler, val, log)

	// ========== OFFER BOUNDED CONTEXT ========== 

	// Offer repositories
	offerRepo := offerPersistence.NewPostgresOfferRepository(db)
	offerCodeRepo := offerPersistence.NewPostgresOfferCodeRepository(db)
	offerItemCriteriaRepo := offerPersistence.NewPostgresOfferItemCriteriaRepository(db)
	offerRuleRepo := offerPersistence.NewPostgresOfferRuleRepository(db)
	offerPriceDataRepo := offerPersistence.NewPostgresOfferPriceDataRepository(db)
	qualCritOfferXrefRepo := offerPersistence.NewPostgresQualCritOfferXrefRepository(db)
	tarCritOfferXrefRepo := offerPersistence.NewPostgresTarCritOfferXrefRepository(db)

	// Offer application services
	offerService := offerApp.NewOfferService(
		offerRepo,
		offerCodeRepo,
		offerItemCriteriaRepo,
		offerRuleRepo,
		offerPriceDataRepo,
		qualCritOfferXrefRepo,
		tarCritOfferXrefRepo,
	)

	// ========== INVENTORY BOUNDED CONTEXT ========== 

	// Inventory repositories
	inventoryLevelRepo := inventoryPersistence.NewPostgresInventoryRepository(db)

	// Inventory application services
	inventoryService := inventoryApp.NewInventoryService(inventoryLevelRepo) // NewInventoryService takes a repo

	// ========== TAX BOUNDED CONTEXT ==========

	// Tax repositories
	// taxDetailRepo := taxPersistence.NewPostgresTaxDetailRepository(db) // Commented out - old tax implementation

	// Tax application services
	// taxService := taxApp.NewTaxService(taxDetailRepo) // Commented out - old tax implementation

	// ========== ORDER BOUNDED CONTEXT ========== 

	// Order repositories
	orderRepo := orderPersistence.NewPostgresOrderRepository(db)
	orderItemRepo := orderPersistence.NewPostgresOrderItemRepository(db)
	orderAdjustmentRepo := orderPersistence.NewPostgresOrderAdjustmentRepository(db)
	orderItemAdjustmentRepo := orderPersistence.NewPostgresOrderItemAdjustmentRepository(db)
	orderItemAttributeRepo := orderPersistence.NewPostgresOrderItemAttributeRepository(db)
	fulfillmentGroupRepo := orderPersistence.NewPostgresFulfillmentGroupRepository(db)

	// Order application service
	orderService := orderApp.NewOrderService(
		orderRepo,
		orderItemRepo,
		orderAdjustmentRepo,
		orderItemAdjustmentRepo,
		orderItemAttributeRepo,
		fulfillmentGroupRepo,
		offerService,
		inventoryService,
		productService,
		skuService,
		// taxService, // Commented out - old tax implementation
	)

	// Order query handlers
	orderQueryHandler := orderQueries.NewOrderQueryHandler(orderService, cacheStore, log)

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
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(middleware.CORSConfig{ // Convert config.CORSConfig to middleware.CORSConfig
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		ExposedHeaders:   cfg.CORS.ExposedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}))

	// Static file serving
	r.Handle("/static/storefront/*", http.StripPrefix("/static/storefront/", http.FileServer(http.Dir("web/storefront/static"))))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Storefront HTML Pages
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := renderer.MergeData(renderer.BaseData(), map[string]interface{}{
			"PageTitle":    "Home",
			"HeroTitle":    "Welcome to Broadleaf Commerce",
			"HeroSubtitle": "Discover amazing products at great prices",
			"FeaturedCategories": []map[string]interface{}{
				{"Name": "Electronics", "URL": "electronics", "ImageURL": "/static/storefront/img/categories/electronics.jpg"},
				{"Name": "Clothing", "URL": "clothing", "ImageURL": "/static/storefront/img/categories/clothing.jpg"},
				{"Name": "Home & Garden", "URL": "home-garden", "ImageURL": "/static/storefront/img/categories/home.jpg"},
			},
			"FeaturedProducts": []map[string]interface{}{},
			"NewArrivals":      []map[string]interface{}{},
		})
		tmplRenderer.RenderHTML(w, "storefront/home", data)
	})

	r.Get("/checkout", func(w http.ResponseWriter, r *http.Request) {
		data := renderer.MergeData(renderer.BaseData(), map[string]interface{}{
			"PageTitle":   "Checkout",
			"CurrentStep": 1,
			"Cart": map[string]interface{}{
				"Subtotal":     "99.99",
				"ShippingCost": "9.99",
				"Tax":          "8.00",
				"Total":        "117.98",
			},
			"ShippingMethods": []map[string]interface{}{
				{"ID": "standard", "Name": "Standard Shipping", "Description": "5-7 business days", "EstimatedDays": "5-7", "Cost": 9.99},
				{"ID": "express", "Name": "Express Shipping", "Description": "2-3 business days", "EstimatedDays": "2-3", "Cost": 19.99},
			},
		})
		tmplRenderer.RenderHTML(w, "storefront/checkout", data)
	})

	// API info
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
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

	log.WithField("contexts", "catalog, customer, order, fulfillment").Info("All storefront contexts initialized")

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Database.MaxIdleTime, // Use a relevant idle timeout from config
	}

	// Start server in a goroutine
	go func() {
		log.WithField("address", addr).Info("Storefront API server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server failed to start")
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
		log.WithError(err).Error("Server forced to shutdown")
	}

	log.Info("Storefront API server stopped")
}

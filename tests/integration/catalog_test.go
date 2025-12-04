package integration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/qhato/ecommerce/internal/catalog/application"
	"github.com/qhato/ecommerce/internal/catalog/domain"
	"github.com/qhato/ecommerce/internal/catalog/infrastructure/persistence"
	"github.com/qhato/ecommerce/pkg/cache"
	"github.com/qhato/ecommerce/pkg/logger"
	"github.com/qhato/ecommerce/pkg/testutil"
)

const productTableSchema = `
CREATE TABLE IF NOT EXISTS blc_product (
    product_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    long_description TEXT,
    product_number VARCHAR(100),
    manufacturer VARCHAR(100),
    is_featured BOOLEAN DEFAULT FALSE,
    can_sell_without_options BOOLEAN DEFAULT TRUE,
    url VARCHAR(255),
    display_template VARCHAR(255),
    active_start_date TIMESTAMP,
    active_end_date TIMESTAMP,
    archived_flag BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

func TestCatalogIntegration_CreateAndRetrieveProduct(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)

	// Create schema
	testDB.CreateTestSchema(t, productTableSchema)

	// Setup dependencies
	log := logger.New("test", "debug")
	productRepo := persistence.NewPostgresProductRepository(testDB.DB)
	cacheImpl := cache.NewInMemoryCache()
	mockBus := testutil.NewMockEventBus()

	// Create command handler
	commandHandler := application.NewProductCommandHandler(productRepo, mockBus, log)

	// Create product
	ctx := context.Background()
	cmd := &application.CreateProductCommand{
		Name:                  "Integration Test Product",
		Description:           "Test Description",
		CanSellWithoutOptions: true,
		Manufacturer:          "Test Manufacturer",
	}

	productDTO, err := commandHandler.CreateProduct(ctx, cmd)
	testutil.AssertNoError(t, err, "CreateProduct should succeed")
	testutil.AssertNotNil(t, productDTO, "ProductDTO should not be nil")
	testutil.AssertEqual(t, productDTO.Name, cmd.Name, "Product name")

	// Retrieve product
	queryHandler := application.NewProductQueryHandler(productRepo, cacheImpl, log)
	retrieved, err := queryHandler.GetProductByID(ctx, productDTO.ID)
	testutil.AssertNoError(t, err, "GetProductByID should succeed")
	testutil.AssertEqual(t, retrieved.ID, productDTO.ID, "Product ID should match")
	testutil.AssertEqual(t, retrieved.Name, productDTO.Name, "Product name should match")

	// Verify event was published
	events := mockBus.GetEventsByType("ProductCreated")
	testutil.AssertLen(t, events, 1, "Should publish ProductCreated event")
}

func TestCatalogIntegration_UpdateProduct(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.CreateTestSchema(t, productTableSchema)

	log := logger.New("test", "debug")
	productRepo := persistence.NewPostgresProductRepository(testDB.DB)
	mockBus := testutil.NewMockEventBus()
	commandHandler := application.NewProductCommandHandler(productRepo, mockBus, log)

	ctx := context.Background()

	// Create product
	createCmd := &application.CreateProductCommand{
		Name:        "Original Name",
		Description: "Original Description",
	}
	product, err := commandHandler.CreateProduct(ctx, createCmd)
	testutil.AssertNoError(t, err, "CreateProduct should succeed")

	// Update product
	updateCmd := &application.UpdateProductCommand{
		ID:          product.ID,
		Name:        "Updated Name",
		Description: "Updated Description",
	}
	updated, err := commandHandler.UpdateProduct(ctx, updateCmd)
	testutil.AssertNoError(t, err, "UpdateProduct should succeed")
	testutil.AssertEqual(t, updated.Name, "Updated Name", "Name should be updated")
	testutil.AssertEqual(t, updated.Description, "Updated Description", "Description should be updated")

	// Verify events
	events := mockBus.GetEvents()
	testutil.AssertTrue(t, len(events) >= 2, "Should have at least 2 events")
}

func TestCatalogIntegration_ProductLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.CreateTestSchema(t, productTableSchema)

	log := logger.New("test", "debug")
	productRepo := persistence.NewPostgresProductRepository(testDB.DB)
	cacheImpl := cache.NewInMemoryCache()
	mockBus := testutil.NewMockEventBus()

	commandHandler := application.NewProductCommandHandler(productRepo, mockBus, log)
	queryHandler := application.NewProductQueryHandler(productRepo, cacheImpl, log)

	ctx := context.Background()

	// 1. Create product
	createCmd := &application.CreateProductCommand{
		Name:                  "Lifecycle Test Product",
		Description:           "Testing full lifecycle",
		CanSellWithoutOptions: true,
	}
	product, err := commandHandler.CreateProduct(ctx, createCmd)
	testutil.AssertNoError(t, err, "Create should succeed")

	// 2. Verify it's retrievable
	retrieved, err := queryHandler.GetProductByID(ctx, product.ID)
	testutil.AssertNoError(t, err, "Retrieve should succeed")
	testutil.AssertFalse(t, retrieved.ArchivedFlag, "Should not be archived")

	// 3. Archive product
	archiveCmd := &application.ArchiveProductCommand{
		ID: product.ID,
	}
	err = commandHandler.ArchiveProduct(ctx, archiveCmd)
	testutil.AssertNoError(t, err, "Archive should succeed")

	// 4. Verify archived
	archived, err := queryHandler.GetProductByID(ctx, product.ID)
	testutil.AssertNoError(t, err, "Should still retrieve archived product")
	testutil.AssertTrue(t, archived.ArchivedFlag, "Should be archived")

	// 5. Delete product
	deleteCmd := &application.DeleteProductCommand{
		ID: product.ID,
	}
	err = commandHandler.DeleteProduct(ctx, deleteCmd)
	testutil.AssertNoError(t, err, "Delete should succeed")

	// 6. Verify deleted
	_, err = queryHandler.GetProductByID(ctx, product.ID)
	testutil.AssertError(t, err, "Should not find deleted product")
}

func TestCatalogIntegration_ListProductsWithPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.CreateTestSchema(t, productTableSchema)

	log := logger.New("test", "debug")
	productRepo := persistence.NewPostgresProductRepository(testDB.DB)
	cacheImpl := cache.NewInMemoryCache()
	mockBus := testutil.NewMockEventBus()

	commandHandler := application.NewProductCommandHandler(productRepo, mockBus, log)
	queryHandler := application.NewProductQueryHandler(productRepo, cacheImpl, log)

	ctx := context.Background()

	// Create 15 products
	for i := 1; i <= 15; i++ {
		cmd := &application.CreateProductCommand{
			Name:        testutil.Sprintf("Product %d", i),
			Description: testutil.Sprintf("Description %d", i),
		}
		_, err := commandHandler.CreateProduct(ctx, cmd)
		testutil.AssertNoError(t, err, "CreateProduct should succeed")
	}

	// List first page
	query := &application.ListProductsQuery{
		Page:     1,
		PageSize: 10,
	}
	results, err := queryHandler.ListProducts(ctx, query)
	testutil.AssertNoError(t, err, "ListProducts should succeed")
	testutil.AssertLen(t, results, 10, "Should return 10 products")

	// List second page
	query.Page = 2
	results, err = queryHandler.ListProducts(ctx, query)
	testutil.AssertNoError(t, err, "ListProducts page 2 should succeed")
	testutil.AssertLen(t, results, 5, "Should return remaining 5 products")
}

func TestCatalogIntegration_ConcurrentProductCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Teardown(t)
	testDB.CreateTestSchema(t, productTableSchema)

	log := logger.New("test", "debug")
	productRepo := persistence.NewPostgresProductRepository(testDB.DB)
	mockBus := testutil.NewMockEventBus()
	commandHandler := application.NewProductCommandHandler(productRepo, mockBus, log)

	ctx := context.Background()
	numGoroutines := 10
	errors := make(chan error, numGoroutines)

	// Create products concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			cmd := &application.CreateProductCommand{
				Name:        testutil.Sprintf("Concurrent Product %d", index),
				Description: "Concurrent test",
			}
			_, err := commandHandler.CreateProduct(ctx, cmd)
			errors <- err
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		testutil.AssertNoError(t, err, "Concurrent create should succeed")
	}

	// Verify all products were created
	cacheImpl := cache.NewInMemoryCache()
	queryHandler := application.NewProductQueryHandler(productRepo, cacheImpl, log)
	results, err := queryHandler.ListProducts(ctx, &application.ListProductsQuery{
		Page:     1,
		PageSize: 100,
	})
	testutil.AssertNoError(t, err, "List should succeed")
	testutil.AssertTrue(t, len(results) >= numGoroutines, "Should have created all products")
}

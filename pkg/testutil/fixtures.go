package testutil

import (
	"time"
)

// Common test fixtures and factory functions

// FixtureProduct returns a sample product for testing
func FixtureProduct() map[string]interface{} {
	return map[string]interface{}{
		"id":                1,
		"name":              "Test Product",
		"description":       "Test product description",
		"long_description":  "Long test product description",
		"product_number":    "PROD-001",
		"manufacturer":      "Test Manufacturer",
		"is_featured":       false,
		"can_sell_without_options": true,
		"url":               "/test-product",
		"display_template":  "default",
		"created_at":        time.Now(),
		"updated_at":        time.Now(),
	}
}

// FixtureCategory returns a sample category for testing
func FixtureCategory() map[string]interface{} {
	return map[string]interface{}{
		"id":               1,
		"name":             "Test Category",
		"description":      "Test category description",
		"long_description": "Long test category description",
		"url":              "/test-category",
		"url_key":          "test-category",
		"active_start_date": time.Now(),
		"created_at":       time.Now(),
		"updated_at":       time.Now(),
	}
}

// FixtureSKU returns a sample SKU for testing
func FixtureSKU(productID int64) map[string]interface{} {
	return map[string]interface{}{
		"id":                1,
		"product_id":        productID,
		"name":              "Test SKU",
		"description":       "Test SKU description",
		"long_description":  "Long test SKU description",
		"sku_number":        "SKU-001",
		"sale_price":        99.99,
		"retail_price":      129.99,
		"cost_price":        50.00,
		"active":            true,
		"active_start_date": time.Now(),
		"created_at":        time.Now(),
		"updated_at":        time.Now(),
	}
}

// FixtureCustomer returns a sample customer for testing
func FixtureCustomer() map[string]interface{} {
	return map[string]interface{}{
		"id":          1,
		"username":    "testuser",
		"email":       "test@example.com",
		"first_name":  "John",
		"last_name":   "Doe",
		"password":    "$2a$10$abcdefghijklmnopqrstuv", // bcrypt hash
		"is_active":   true,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
}

// FixtureOrder returns a sample order for testing
func FixtureOrder(customerID int64) map[string]interface{} {
	return map[string]interface{}{
		"id":           1,
		"customer_id":  customerID,
		"order_number": "ORD-001",
		"status":       "PENDING",
		"total":        199.99,
		"subtotal":     199.99,
		"tax_total":    0.00,
		"shipping_total": 0.00,
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}
}

// FixtureOrderItem returns a sample order item for testing
func FixtureOrderItem(orderID, skuID int64) map[string]interface{} {
	return map[string]interface{}{
		"id":         1,
		"order_id":   orderID,
		"sku_id":     skuID,
		"name":       "Test Product",
		"quantity":   2,
		"price":      99.99,
		"total":      199.98,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
}

// FixturePayment returns a sample payment for testing
func FixturePayment(orderID, customerID int64) map[string]interface{} {
	return map[string]interface{}{
		"id":             1,
		"order_id":       orderID,
		"customer_id":    customerID,
		"amount":         199.99,
		"payment_type":   "CREDIT_CARD",
		"status":         "PENDING",
		"transaction_id": "TXN-001",
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}
}

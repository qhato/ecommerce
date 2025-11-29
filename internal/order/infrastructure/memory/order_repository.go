package memory

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/qhato/ecommerce/internal/order/domain"
)

// OrderRepository implements domain.OrderRepository for in-memory persistence.
type OrderRepository struct {
	mu     sync.RWMutex
	orders map[int64]*domain.Order
	nextID int64
}

// NewOrderRepository creates a new in-memory order repository.
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int64]*domain.Order),
		nextID: 1,
	}
}

// Create stores a new order.
func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = r.nextID
	order.OrderNumber = generateOrderNumber(order.ID)
	r.nextID++

	r.orders[order.ID] = order
	return nil
}

// Update updates an existing order.
func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID]; !ok {
		return fmt.Errorf("order with ID %d not found for update", order.ID)
	}
	r.orders[order.ID] = order
	return nil
}

// FindByID retrieves an order by its unique identifier.
func (r *OrderRepository) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, nil
	}
	return order, nil
}

// FindByOrderNumber retrieves an order by its order number.
func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, order := range r.orders {
		if order.OrderNumber == orderNumber {
			return order, nil
		}
	}
	return nil, nil
}

// FindByCustomerID retrieves orders by customer ID.
func (r *OrderRepository) FindByCustomerID(ctx context.Context, customerID int64, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var customerOrders []*domain.Order
	for _, order := range r.orders {
		if order.CustomerID == customerID {
			customerOrders = append(customerOrders, order)
		}
	}
	// TODO: Implement filtering, pagination, and sorting based on the filter
	return customerOrders, int64(len(customerOrders)), nil
}

// FindAll retrieves all orders.
func (r *OrderRepository) FindAll(ctx context.Context, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	// TODO: Implement filtering, pagination, and sorting based on the filter
	return orders, int64(len(orders)), nil
}

func generateOrderNumber(orderID int64) string {
	// Simple generation for in-memory, can be more complex with prefix/suffix/timestamp
	seed := time.Now().UnixNano() + int64(rand.Intn(1000))
	rand.Seed(seed)
	return fmt.Sprintf("ORD-%s-%d", strconv.FormatInt(time.Now().Unix(), 36), rand.Intn(10000))
}

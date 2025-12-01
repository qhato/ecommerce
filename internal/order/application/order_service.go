package application

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	catalogApp "github.com/qhato/ecommerce/internal/catalog/application"
	inventoryApp "github.com/qhato/ecommerce/internal/inventory/application"
	offerApp "github.com/qhato/ecommerce/internal/offer/application"
	offerDomain "github.com/qhato/ecommerce/internal/offer/domain"
	"github.com/qhato/ecommerce/internal/order/domain"
	taxApp "github.com/qhato/ecommerce/internal/tax/application"
)

// OrderService defines the application service for order-related operations.
type OrderService interface {
	// CreateOrder creates a new order.
	CreateOrder(ctx context.Context, cmd *CreateOrderCommand) (*OrderDTO, error)

	// HandleGetOrderByID retrieves an order by its ID.
	HandleGetOrderByID(ctx context.Context, id int64) (*OrderDTO, error)

	// UpdateOrderStatus updates the status of an existing order.
	UpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error

	// AddItemToOrder adds an item to an existing order.
	AddItemToOrder(ctx context.Context, orderID int64, cmd *AddItemToOrderCommand) (*OrderItemDTO, error)

	// UpdateOrderItemQuantity updates the quantity of an existing order item.
	UpdateOrderItemQuantity(ctx context.Context, orderItemID int64, newQuantity int) (*OrderItemDTO, error)

	// RemoveOrderItem removes an item from the order.
	RemoveOrderItem(ctx context.Context, orderItemID int64) error

	// SubmitOrder submits an order for processing.
	SubmitOrder(ctx context.Context, orderID int64) error

	// CancelOrder cancels an existing order.
	CancelOrder(ctx context.Context, orderID int64, reason string) error

	// ApplyOffersToOrder fetches active offers and applies them to an order.
	ApplyOffersToOrder(ctx context.Context, orderID int64, customerID int64, couponCode *string) (*OrderDTO, error)

	// CreateFulfillmentGroup creates a new fulfillment group for an order.
	CreateFulfillmentGroup(ctx context.Context, orderID int64, cmd *CreateFulfillmentGroupCommand) (*FulfillmentGroupDTO, error)

	// UpdateOrderShippingDetails updates the shipping details of an order.
	UpdateOrderShippingDetails(ctx context.Context, orderID int64, shippingCost float64) error

	// ListOrders retrieves a paginated list of orders based on a filter.
	ListOrders(ctx context.Context, filter *domain.OrderFilter) ([]*domain.Order, int64, error)

	// GetOrderByOrderNumber retrieves an order by its order number.
	GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*OrderDTO, error)
}

// CreateOrderCommand is a command to create a new order.
type CreateOrderCommand struct {
	CustomerID   int64
	EmailAddress string
	Name         string
	CurrencyCode string
	LocaleCode   string
	IsPreview    bool
	TaxOverride  bool
}

// AddItemToOrderCommand is a command to add an item to an order.
type AddItemToOrderCommand struct {
	SKUID        int64
	Quantity     int
	TaxCategory  string
	CategoryID   *int64
	GiftWrapItemID *int64
	ParentOrderItemID *int64
	PersonalMessageID *int64
	// Additional fields for OrderItem creation can be added here.
}

// CreateFulfillmentGroupCommand is a command to create a new fulfillment group.
type CreateFulfillmentGroupCommand struct {
	Type        string
	AddressID   *int64
	FulfillmentOptionID *int64
	PersonalMessageID *int64
	PhoneID   *int64
	IsPrimary bool
	Status    string
	// Other fields for fulfillment group
}

type orderService struct {
	orderRepo               domain.OrderRepository
	orderItemRepo           domain.OrderItemRepository
	orderAdjustmentRepo     domain.OrderAdjustmentRepository
	orderItemAdjustmentRepo domain.OrderItemAdjustmentRepository
	orderItemAttributeRepo  domain.OrderItemAttributeRepository
	fulfillmentGroupRepo    domain.FulfillmentGroupRepository
	offerService            offerApp.OfferService
	inventoryService        inventoryApp.InventoryService
	productService          catalogApp.ProductService
	skuService              catalogApp.SkuService
	taxService              taxApp.TaxService
}

// NewOrderService creates a new instance of OrderService.
func NewOrderService(
	orderRepo domain.OrderRepository,
	orderItemRepo domain.OrderItemRepository,
	orderAdjustmentRepo domain.OrderAdjustmentRepository,
	orderItemAdjustmentRepo domain.OrderItemAdjustmentRepository,
	orderItemAttributeRepo domain.OrderItemAttributeRepository,
	fulfillmentGroupRepo domain.FulfillmentGroupRepository,
	offerService offerApp.OfferService,
	inventoryService inventoryApp.InventoryService,
	productService catalogApp.ProductService,
	skuService catalogApp.SkuService,
	taxService taxApp.TaxService,
) OrderService {
	return &orderService{
		orderRepo:               orderRepo,
		orderItemRepo:           orderItemRepo,
		orderAdjustmentRepo:     orderAdjustmentRepo,
		orderItemAdjustmentRepo: orderItemAdjustmentRepo,
		orderItemAttributeRepo:  orderItemAttributeRepo,
		fulfillmentGroupRepo:    fulfillmentGroupRepo,
		offerService:            offerService,
		inventoryService:        inventoryService,
		productService:          productService,
		skuService:              skuService,
		taxService:              taxService,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, cmd *CreateOrderCommand) (*OrderDTO, error) {
	order := domain.NewOrder(cmd.CustomerID, cmd.EmailAddress, cmd.Name, cmd.CurrencyCode, cmd.LocaleCode)
	order.IsPreview = cmd.IsPreview
	order.TaxOverride = cmd.TaxOverride

	err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return ToOrderDTO(order), nil
}

// HandleGetOrderByID handles the get order by ID query
func (s *orderService) HandleGetOrderByID(ctx context.Context, id int64) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order by ID: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order with ID %d not found", id)
	}

	items, err := s.orderItemRepo.FindByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items for order %d: %w", id, err)
	}
	orderAdjustments, err := s.orderAdjustmentRepo.FindByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order adjustments for order %d: %w", id, err)
	}
	fulfillmentGroups, err := s.fulfillmentGroupRepo.FindByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fulfillment groups for order %d: %w", id, err)
	}

	return toOrderDTOWithRelations(order, items, orderAdjustments, fulfillmentGroups), nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order by ID for status update: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order with ID %d not found for status update", orderID)
	}

	order.UpdateStatus(status)
	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

func (s *orderService) AddItemToOrder(ctx context.Context, orderID int64, cmd *AddItemToOrderCommand) (*OrderItemDTO, error) {
	// 1. Get SKU details
	skuDTO, err := s.skuService.GetSkuByID(ctx, cmd.SKUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get SKU details for ID %d: %w", cmd.SKUID, err)
	}
	if skuDTO == nil {
		return nil, fmt.Errorf("SKU with ID %d not found", cmd.SKUID)
	}

	// 2. Get Product details from SKU's DefaultProductID
	var productID int64
	if skuDTO.DefaultProductID != nil {
		productID = *skuDTO.DefaultProductID
	} else {
		return nil, fmt.Errorf("SKU with ID %d has no associated default product", cmd.SKUID)
	}

	// 3. Allocate inventory
	skuAvailability, err := s.inventoryService.GetInventoryLevelBySKUID(ctx, strconv.FormatInt(cmd.SKUID, 10)) // Use new method
	if err != nil || skuAvailability == nil {
		return nil, fmt.Errorf("failed to get SKU availability for ID %d: %w", cmd.SKUID, err)
	}
	if skuAvailability.QuantityOnHand < cmd.Quantity { // Use QuantityOnHand
		return nil, fmt.Errorf("not enough quantity on hand for SKU %d", cmd.SKUID)
	}

	// Update inventory quantities
	updatedLevel, err := s.inventoryService.UpdateInventoryQuantities(
		ctx,
		skuAvailability.ID,
		skuAvailability.QuantityOnHand-cmd.Quantity,
		skuAvailability.QuantityReserved+cmd.Quantity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate inventory for SKU %d: %w", cmd.SKUID, err)
	}

	// 4. Create OrderItem domain entity
	item, err := domain.NewOrderItem(
		orderID,
		cmd.SKUID,
		productID,
		skuDTO.Name, // Use SKU name as item name
		cmd.Quantity,
		skuDTO.RetailPrice,
		skuDTO.SalePrice,
		cmd.TaxCategory,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order item domain entity: %w", err)
	}

	item.CategoryID = cmd.CategoryID
	item.GiftWrapItemID = cmd.GiftWrapItemID
	item.ParentOrderItemID = cmd.ParentOrderItemID
	item.PersonalMessageID = cmd.PersonalMessageID

	// Calculate initial tax based on TaxService (simplified)
	taxAmount := 0.0
	if cmd.TaxCategory != "" {
		taxAmount, err = s.taxService.CalculateTaxForItem(ctx, orderID, item.TotalPrice, cmd.TaxCategory)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax for item: %w", err)
		}
	}
	item.SetTaxAmount(taxAmount)

	// 5. Save OrderItem
	err = s.orderItemRepo.Save(ctx, item)
	if err != nil {
		// Attempt to deallocate inventory if item save fails
		_, deallocErr := s.inventoryService.UpdateInventoryQuantities(
			ctx,
			updatedLevel.ID, // Use the ID from the updated level
			updatedLevel.QuantityOnHand+cmd.Quantity,
			updatedLevel.QuantityReserved-cmd.Quantity,
		)
		if deallocErr != nil {
			return nil, fmt.Errorf("failed to save order item: %w (and failed to deallocate inventory: %v)", err, deallocErr)
		}
		return nil, fmt.Errorf("failed to save order item: %w", err)
	}

	// 6. Recalculate order totals
	// The order totals will be recalculated by ApplyOffersToOrder or a dedicated recalculate method
	// For now, we update the order's top-level totals after each item add/update/remove
	order, err := s.orderRepo.FindByID(ctx, orderID) // Re-fetch order to ensure consistency
	if err != nil {
		return nil, fmt.Errorf("failed to re-fetch order to recalculate totals: %w", err)
	}
	order.OrderSubtotal += item.TotalPrice
	order.TotalTax += item.TaxAmount
	order.OrderTotal = order.OrderSubtotal + order.TotalTax + order.TotalShipping // Assuming shipping is calculated elsewhere

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order totals: %w", err)
	}

	return ToOrderItemDTO(item), nil
}

func (s *orderService) UpdateOrderItemQuantity(ctx context.Context, orderItemID int64, newQuantity int) (*OrderItemDTO, error) {
	item, err := s.orderItemRepo.FindByID(ctx, orderItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order item by ID: %w", err)
	}
	if item == nil {
		return nil, fmt.Errorf("order item with ID %d not found", orderItemID)
	}

	order, err := s.orderRepo.FindByID(ctx, item.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order by ID for item update: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order with ID %d not found for item update", item.OrderID)
	}

	oldQuantity := item.Quantity
	quantityDiff := newQuantity - oldQuantity

	if quantityDiff != 0 {
		skuAvailability, err := s.inventoryService.GetInventoryLevelBySKUID(ctx, strconv.FormatInt(item.SKUID, 10))
		if err != nil || skuAvailability == nil {
			return nil, fmt.Errorf("failed to get SKU availability for ID %d: %w", item.SKUID, err)
		}

		if quantityDiff > 0 { // Increasing quantity, need to allocate more
			if skuAvailability.QuantityOnHand < quantityDiff {
				return nil, fmt.Errorf("not enough quantity on hand for SKU %d to increase by %d", item.SKUID, quantityDiff)
			}
			_, err = s.inventoryService.UpdateInventoryQuantities(
				ctx,
				skuAvailability.ID,
				skuAvailability.QuantityOnHand-quantityDiff,
				skuAvailability.QuantityReserved+quantityDiff,
			)
		} else { // Decreasing quantity, need to deallocate
			_, err = s.inventoryService.UpdateInventoryQuantities(
				ctx,
				skuAvailability.ID,
				skuAvailability.QuantityOnHand+quantityDiff, // This should be QuantityOnHand + abs(quantityDiff)
				skuAvailability.QuantityReserved-quantityDiff, // This should be QuantityReserved - abs(quantityDiff)
			)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to adjust inventory for SKU %d: %w", item.SKUID, err)
		}
	}

	err = item.UpdateQuantity(newQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to update order item quantity: %w", err)
	}

	// Recalculate tax for the item
	taxAmount := 0.0
	if item.TaxCategory != "" {
		taxAmount, err = s.taxService.CalculateTaxForItem(ctx, order.ID, item.TotalPrice, item.TaxCategory)
		if err != nil {
			return nil, fmt.Errorf("failed to recalculate tax for item: %w", err)
		}
	}
	item.SetTaxAmount(taxAmount)

	err = s.orderItemRepo.Save(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to save order item after quantity update: %w", err)
	}

	// Recalculate order totals
	order.OrderSubtotal += (item.TotalPrice - (item.Price * float64(oldQuantity))) // Adjust subtotal by change
	order.TotalTax += (item.TaxAmount - (taxAmount * float64(oldQuantity)))         // Adjust total tax
	order.OrderTotal = order.OrderSubtotal + order.TotalTax + order.TotalShipping

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order totals after item quantity update: %w", err)
	}

	return ToOrderItemDTO(item), nil
}

func (s *orderService) RemoveOrderItem(ctx context.Context, orderItemID int64) error {
	item, err := s.orderItemRepo.FindByID(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to find order item by ID: %w", err)
	}
	if item == nil {
		return fmt.Errorf("order item with ID %d not found", orderItemID)
	}

	order, err := s.orderRepo.FindByID(ctx, item.OrderID)
	if err != nil {
		return fmt.Errorf("failed to find order by ID for item removal: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order with ID %d not found for item removal", item.OrderID)
	}

	// Deallocate inventory
	skuAvailability, err := s.inventoryService.GetInventoryLevelBySKUID(ctx, strconv.FormatInt(item.SKUID, 10))
	if err != nil || skuAvailability == nil {
		return fmt.Errorf("failed to get SKU availability for ID %d: %w", item.SKUID, err)
	}
	_, err = s.inventoryService.UpdateInventoryQuantities(
		ctx,
		skuAvailability.ID,
		skuAvailability.QuantityOnHand+item.Quantity,
		skuAvailability.QuantityReserved-item.Quantity,
	)
	if err != nil {
		return fmt.Errorf("failed to deallocate inventory for SKU %d: %w", item.SKUID, err)
	}

	// Delete item and associated entities
	err = s.orderItemAttributeRepo.DeleteByOrderItemID(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item attributes for item %d: %w", orderItemID, err)
	}
	err = s.orderItemAdjustmentRepo.DeleteByOrderItemID(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item adjustments for item %d: %w", orderItemID, err)
	}
	err = s.orderItemRepo.Delete(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	// Recalculate order totals
	order.OrderSubtotal -= item.TotalPrice
	order.TotalTax -= item.TaxAmount
	order.OrderTotal = order.OrderSubtotal + order.TotalTax + order.TotalShipping

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to update order totals after item removal: %w", err)
	}

	return nil
}

func (s *orderService) SubmitOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order by ID for submission: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order with ID %d not found for submission", orderID)
	}

	// In a real system, would check if items exist here. Assume application layer handles this.
	// We also assume tax calculation is final before submission.

	err = order.Submit()
	if err != nil {
		return fmt.Errorf("failed to submit order: %w", err)
	}

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to update order after submission: %w", err)
	}
	return nil
}

func (s *orderService) CancelOrder(ctx context.Context, orderID int64, reason string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order by ID for cancellation: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order with ID %d not found for cancellation", orderID)
	}

	if !order.IsCancellable() {
		return fmt.Errorf("order with ID %d is not cancellable in status %s", orderID, order.Status)
	}

	// Deallocate inventory for all items in the order
	items, err := s.orderItemRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order items for deallocation: %w", err)
	}

	for _, item := range items {
		skuAvailability, err := s.inventoryService.GetInventoryLevelBySKUID(ctx, strconv.FormatInt(item.SKUID, 10))
		if err != nil || skuAvailability == nil {
			fmt.Printf("warning: failed to get SKU availability for SKU %d (order %d): %v\n", item.SKUID, orderID, err)
			continue
		}
		_, deallocErr := s.inventoryService.UpdateInventoryQuantities(
			ctx,
			skuAvailability.ID,
			skuAvailability.QuantityOnHand+item.Quantity,
			skuAvailability.QuantityReserved-item.Quantity,
		)
		if deallocErr != nil {
			// Log the error but continue with order cancellation to avoid blocking
			fmt.Printf("warning: failed to deallocate inventory for SKU %d (order %d): %v\n", item.SKUID, orderID, deallocErr)
		}
	}

	order.Cancel()
	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

func (s *orderService) ApplyOffersToOrder(ctx context.Context, orderID int64, customerID int64, couponCode *string) (*OrderDTO, error) {
	// Load the full order graph
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order by ID: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order with ID %d not found", orderID)
	}

	items, err := s.orderItemRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items for order %d: %w", orderID, err)
	}

	// Clear existing adjustments before reapplying
	err = s.orderAdjustmentRepo.DeleteByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear existing order adjustments: %w", err)
	}
	for _, item := range items {
		err = s.orderItemAdjustmentRepo.DeleteByOrderItemID(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to clear existing item adjustments for item %d: %w", item.ID, err)
		}
		// Reset item prices for recalculation
		item.UpdatePrices(item.RetailPrice, item.SalePrice, item.RetailPrice) // Use original retail for base
		err = s.orderItemRepo.Save(ctx, item)
		if err != nil {
			return nil, fmt.Errorf("failed to reset item prices for item %d: %w", item.ID, err)
		}
	}

	// 1. Get all active offers
	activeOffersDTO, err := s.offerService.GetActiveOffers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active offers: %w", err)
	}
	var activeOffers []*offerDomain.Offer
	for _, dto := range activeOffersDTO {
		activeOffers = append(activeOffers, offerApp.ToOfferDomain(*dto)) // Corrected: use offerApp.ToOfferDomain and dereference dto
	}


	var applicableOffers []*offerDomain.Offer

	// Add offers by coupon code if provided
	if couponCode != nil && *couponCode != "" {
		couponOfferDTO, err := s.offerService.GetOfferByCode(ctx, *couponCode)
		if err != nil {
			return nil, fmt.Errorf("failed to find offer by coupon code %s: %w", *couponCode, err)
		}
		if couponOfferDTO != nil && !couponOfferDTO.Archived {
			// Further check customer-specific max uses and audience here if needed
			applicableOffers = append(applicableOffers, offerApp.ToOfferDomain(*couponOfferDTO)) // Corrected
		}
	}

	// Add other automatically applying offers (not requiring a coupon code)
	for _, offer := range activeOffers {
		if offer.AutomaticallyAdded {
			applicableOffers = append(applicableOffers, offer)
		}
	}

	// Sort offers by priority (lower number = higher priority)
	sort.Slice(applicableOffers, func(i, j int) bool {
		return applicableOffers[i].OfferPriority < applicableOffers[j].OfferPriority
	})

	for _, offer := range applicableOffers {
		// Simplified offer application logic. Real logic would be much more complex.
		if offer.OrderMinTotal > 0 && order.OrderSubtotal < offer.OrderMinTotal {
			continue // Order does not meet minimum subtotal
		}

		switch offer.OfferDiscountType {
		case offerDomain.OfferDiscountTypeAmountOff, offerDomain.OfferDiscountTypePercentDiscount:
			if offer.AdjustmentType == offerDomain.OfferAdjustmentTypeOrder {
				// Apply order-level discount
				adjustmentAmount := 0.0
				if offer.OfferType == offerDomain.OfferTypePercentageOff {
					adjustmentAmount = order.OrderSubtotal * offer.OfferValue
				} else if offer.OfferType == offerDomain.OfferTypeAmountOff {
					adjustmentAmount = offer.OfferValue
				}
				if adjustmentAmount > 0 {
					adj, _ := domain.NewOrderAdjustment(order.ID, offer.ID, offer.OfferDescription, -adjustmentAmount, false)
					err = s.orderAdjustmentRepo.Save(ctx, adj)
					if err != nil {
						return nil, fmt.Errorf("failed to save order adjustment: %w", err)
					}
					// Increment offer uses (needs to be handled by offer service)
					// s.offerService.IncrementOfferUses(ctx, offer.ID)
				}
			} else if offer.AdjustmentType == offerDomain.OfferAdjustmentTypeOrderItem {
				// Apply item-level discount
				for _, item := range items {
					// Placeholder for complex item eligibility checks using QualCritOfferXref and TarCritOfferXref
					itemApplies := s.checkItemEligibility(ctx, item, offer)

					if itemApplies {
						itemAdjustmentAmount := 0.0
						basePrice := item.Price // Use current item price, which might already be discounted by higher priority offers
						if offer.OfferType == offerDomain.OfferTypePercentageOff {
							itemAdjustmentAmount = basePrice * offer.OfferValue * float64(item.Quantity)
						} else if offer.OfferType == offerDomain.OfferTypeAmountOff {
							itemAdjustmentAmount = offer.OfferValue * float64(item.Quantity)
						}

							if itemAdjustmentAmount > 0 {
								itemAdj, _ := domain.NewOrderItemAdjustment(item.ID, offer.ID, offer.OfferDescription, -itemAdjustmentAmount, offer.ApplyToSalePrice) // Removed last arg
							err = s.orderItemAdjustmentRepo.Save(ctx, itemAdj)
							if err != nil {
								return nil, fmt.Errorf("failed to save order item adjustment: %w", err)
							}
							
							item.UpdatePrices(item.RetailPrice, item.SalePrice, item.Price-(itemAdjustmentAmount/float64(item.Quantity)))
							err = s.orderItemRepo.Save(ctx, item)
							if err != nil {
								return nil, fmt.Errorf("failed to update order item prices: %w", err)
							}
							// Increment offer uses (needs to be handled by offer service)
							// s.offerService.IncrementOfferUses(ctx, offer.ID)
						}
					}
				}
			}
		// TODO: Implement BOGO logic, Shipping discounts, and more complex rules.
		}
	}

	// Recalculate full order totals after all offers applied
	order.OrderSubtotal = 0.0
	order.TotalTax = 0.0
	order.TotalShipping = 0.0 // Assuming this will be calculated by a shipping service

	for _, item := range items {
		order.OrderSubtotal += item.TotalPrice
		order.TotalTax += item.TaxAmount
	}
	// Sum order adjustments
	orderAdjustments, err := s.orderAdjustmentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order adjustments for total recalculation: %w", err)
	}
	for _, adj := range orderAdjustments {
		order.OrderSubtotal += adj.AdjustmentValue // Adjust subtotal based on order-level discounts
	}
	
	order.OrderTotal = order.OrderSubtotal + order.TotalTax + order.TotalShipping

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order after applying offers: %w", err)
	}

	return toOrderDTOWithRelations(order, items, orderAdjustments, nil), nil // Fulfillment groups not updated here
}

func (s *orderService) CreateFulfillmentGroup(ctx context.Context, orderID int64, cmd *CreateFulfillmentGroupCommand) (*FulfillmentGroupDTO, error) {
	fg, err := domain.NewFulfillmentGroup(orderID, cmd.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to create fulfillment group domain entity: %w", err)
	}

	fg.AddressID = cmd.AddressID
	fg.FulfillmentOptionID = cmd.FulfillmentOptionID
	fg.PersonalMessageID = cmd.PersonalMessageID
	fg.PhoneID = cmd.PhoneID
	fg.IsPrimary = cmd.IsPrimary
	fg.Status = cmd.Status

	err = s.fulfillmentGroupRepo.Save(ctx, fg)
	if err != nil {
		return nil, fmt.Errorf("failed to save fulfillment group: %w", err)
	}

	return ToFulfillmentGroupDTO(fg), nil
}

// UpdateOrderShippingDetails updates the shipping details of an order.
func (s *orderService) UpdateOrderShippingDetails(ctx context.Context, orderID int64, shippingCost float64) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order by ID for shipping update: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order with ID %d not found for shipping update", orderID)
	}

	order.TotalShipping = shippingCost
	order.OrderTotal = order.OrderSubtotal + order.TotalTax + order.TotalShipping // Recalculate total

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to update order shipping details: %w", err)
	}
	return nil
}

// checkItemEligibility is a placeholder for complex item eligibility logic.
// In a real system, this would evaluate offer.OfferItemQualifierRule
// and offer.OfferItemTargetRule against the item's properties.
func (s *orderService) checkItemEligibility(ctx context.Context, item *domain.OrderItem, offer *offerDomain.Offer) bool {
	// For now, a very basic check
	// This should involve:
	// 1. Fetching OfferItemCriteria using offer.OfferItemQualifierRule and offer.OfferItemTargetRule
	// 2. Evaluating match rules against item properties (SKU, Product, Category, quantity, etc.)

	// Placeholder: Assume all items are eligible if no specific rules are set
	return true
}

func toOrderDTO(order *domain.Order) *OrderDTO {
	return &OrderDTO{
		ID:                      order.ID,
		OrderNumber:             order.OrderNumber,
		CustomerID:              order.CustomerID,
		EmailAddress:            order.EmailAddress,
		Name:                    order.Name,
		Status:                  order.Status,
		OrderSubtotal:           order.OrderSubtotal,
		TotalTax:                order.TotalTax,
		TotalShipping:           order.TotalShipping,
		OrderTotal:              order.OrderTotal,
		CurrencyCode:            order.CurrencyCode,
		IsPreview:               order.IsPreview,
		TaxOverride:             order.TaxOverride,
		LocaleCode:              order.LocaleCode,
		SubmitDate:              order.SubmitDate,
		CreatedAt:               order.CreatedAt,
		UpdatedAt:               order.UpdatedAt,
	}
}

func toOrderDTOWithRelations(
	order *domain.Order,
	items []*domain.OrderItem,
	orderAdjustments []*domain.OrderAdjustment,
	fulfillmentGroups []*domain.FulfillmentGroup,
) *OrderDTO {
	orderDTO := toOrderDTO(order)

	itemsDTO := make([]*OrderItemDTO, len(items))
	for i, item := range items {
		itemsDTO[i] = ToOrderItemDTO(item)
	}
	orderDTO.Items = itemsDTO

	adjustmentsDTO := make([]*OrderAdjustmentDTO, len(orderAdjustments))
	for i, adj := range orderAdjustments {
		adjustmentsDTO[i] = ToOrderAdjustmentDTO(adj)
	}
	orderDTO.OrderAdjustments = adjustmentsDTO

	fulfillmentGroupsDTO := make([]*FulfillmentGroupDTO, len(fulfillmentGroups))
	for i, fg := range fulfillmentGroups {
		fulfillmentGroupsDTO[i] = ToFulfillmentGroupDTO(fg)
	}
	orderDTO.FulfillmentGroups = fulfillmentGroupsDTO

	return orderDTO
}

func (s *orderService) ListOrders(ctx context.Context, filter *domain.OrderFilter) ([]*domain.Order, int64, error) {
	orders, total, err := s.orderRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, total, nil
}

func (s *orderService) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to find order by order number %s: %w", orderNumber, err)
	}
	if order == nil {
		return nil, fmt.Errorf("order with order number %s not found", orderNumber)
	}
	return ToOrderDTO(order), nil
}

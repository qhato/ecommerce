package queries

// Inventory Level Queries
type GetInventoryLevelQuery struct {
	ID string
}

type GetInventoryBySKUQuery struct {
	SKUID string
}

type GetInventoryByWarehouseQuery struct {
	WarehouseID string
}

type CheckInventoryAvailabilityQuery struct {
	SKUID    string
	Quantity int
}

type GetLowStockItemsQuery struct {
	WarehouseID *string
	Limit       int
}

type GetBackorderableItemsQuery struct {
	WarehouseID *string
}

// Reservation Queries
type GetReservationQuery struct {
	ID string
}

type GetReservationsByOrderQuery struct {
	OrderID string
}

type GetExpiredReservationsQuery struct{}

type GetActiveReservationsQuery struct {
	SKUID *string
}

// Analytics Queries
type GetInventoryValueQuery struct {
	WarehouseID *string
}

type GetInventoryTurnoverQuery struct {
	SKUID       string
	DaysBack    int
}

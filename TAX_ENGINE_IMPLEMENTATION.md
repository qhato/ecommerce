# Tax Engine Implementation

## Overview

The Tax Engine provides comprehensive tax calculation and management capabilities for the e-commerce platform. It supports complex, multi-jurisdictional tax scenarios including federal, state, county, city, and district-level taxes with support for tax exemptions, compound taxes, and category-based taxation.

## Architecture

The Tax Engine follows hexagonal architecture with clear separation of concerns:

```
internal/tax/
├── domain/              # Business logic and entities
│   ├── tax_jurisdiction.go
│   ├── tax_rate.go
│   ├── tax_calculation.go
│   ├── errors.go
│   ├── events.go
│   └── repository.go
├── application/         # Use cases and orchestration
│   ├── commands/
│   │   ├── tax_commands.go
│   │   └── tax_command_handler.go
│   └── queries/
│       ├── dto.go
│       └── tax_calculator_service.go
├── infrastructure/      # External implementations
│   └── persistence/
│       ├── tax_jurisdiction_repository.go
│       ├── tax_rate_repository.go
│       └── tax_exemption_repository.go
└── ports/              # API interfaces
    └── http/
        └── tax_handler.go
```

## Domain Model

### Core Entities

#### 1. TaxJurisdiction
Represents a tax authority (federal, state, county, city, district):
- **Code**: Unique identifier (e.g., "US-CA", "US-CA-SF")
- **Name**: Human-readable name
- **Type**: FEDERAL, STATE, COUNTY, CITY, DISTRICT
- **Location**: Country, state/province, county, city, postal code
- **Priority**: Application order (lower = applied first)
- **Parent**: Hierarchical relationships

#### 2. TaxRate
Defines tax rates for specific jurisdictions:
- **Types**: PERCENTAGE, FLAT, COMPOUND
- **Categories**: GENERAL, FOOD, CLOTHING, DIGITAL, SHIPPING, SERVICE, EXEMPT
- **Rate**: Tax percentage or flat amount
- **Compound**: Whether calculated on subtotal + existing taxes
- **Thresholds**: Min/max amounts for applicability
- **Shipping Taxable**: Whether applies to shipping charges
- **Date Range**: Start and end dates for validity

#### 3. TaxCalculationRequest
Input for tax calculation:
- **Address**: Shipping and billing addresses
- **Items**: List of taxable items with categories
- **Shipping Amount**: Delivery cost
- **Customer ID**: For exemption lookup
- **Order ID**: For tracking purposes

#### 4. TaxCalculationResult
Output from tax calculation:
- **Items**: Taxed items with detailed breakdowns
- **Total Tax**: Aggregate tax amount
- **Breakdowns**: Tax grouped by jurisdiction
- **Jurisdictions Used**: List of applied jurisdictions
- **Effective Tax Rate**: Overall tax percentage

#### 5. TaxExemption
Customer or certificate-based exemptions:
- **Customer ID**: Who has the exemption
- **Certificate**: Unique certificate number
- **Jurisdiction**: Specific or all jurisdictions
- **Category**: Specific or all categories
- **Date Range**: Active period

## Business Logic

### Tax Calculation Algorithm

1. **Jurisdiction Resolution**
   - Match shipping address to applicable jurisdictions
   - Consider federal, state, county, city, district levels
   - Sort by priority (lower first)

2. **Rate Selection**
   - Find rates for matched jurisdictions
   - Filter by tax category
   - Check date validity
   - Apply thresholds

3. **Exemption Application**
   - Lookup customer exemptions
   - Check certificate validity
   - Match jurisdiction and category
   - Verify date range

4. **Tax Calculation**
   - Calculate simple taxes: subtotal × rate
   - Calculate compound taxes: (subtotal + previous taxes) × rate
   - Apply flat taxes: rate × quantity
   - Sum all applicable taxes

5. **Breakdown Generation**
   - Group taxes by jurisdiction
   - Track individual tax applications
   - Calculate effective rates
   - Provide detailed audit trail

### Compound Tax Support

Compound taxes are calculated on top of other taxes:
```
Item Subtotal: $100
Tax 1 (5%): $100 × 0.05 = $5
Tax 2 (2%, compound): ($100 + $5) × 0.02 = $2.10
Total Tax: $7.10
```

### Category-Based Taxation

Different items can be taxed at different rates:
- **GENERAL**: Standard merchandise (highest rates)
- **FOOD**: Groceries (often reduced or exempt)
- **CLOTHING**: Apparel (may have thresholds)
- **DIGITAL**: Digital goods and services
- **SHIPPING**: Delivery charges
- **SERVICE**: Professional services
- **EXEMPT**: Tax-free items

## Implementation Details

### Domain Layer (6 files, ~850 lines)

**tax_jurisdiction.go** (130 lines)
- 5 jurisdiction types
- Geographic matching logic
- Parent-child relationships
- Activation/deactivation

**tax_rate.go** (180 lines)
- 3 rate calculation types
- 7 tax categories
- Compound tax support
- Threshold validation
- Date range checking

**tax_calculation.go** (250 lines)
- Request/response structures
- Address and item models
- Breakdown structures
- Exemption models

**errors.go** (60 lines)
- 15+ domain-specific errors
- Validation errors
- Business rule violations

**events.go** (160 lines)
- 8 domain events
- Event metadata
- Event type constants

**repository.go** (120 lines)
- 4 repository interfaces
- 30+ query methods
- Bulk operations

### Application Layer (4 files, ~1,400 lines)

**commands/tax_commands.go** (110 lines)
- 9 command DTOs
- Create/update/delete commands
- Bulk operations

**commands/tax_command_handler.go** (400 lines)
- Command processing
- Validation logic
- Repository coordination
- Error handling

**queries/tax_calculator_service.go** (500 lines)
- Core calculation engine
- Jurisdiction matching
- Rate application
- Exemption processing
- Query methods (15+)

**queries/dto.go** (390 lines)
- 10+ DTO types
- Domain-to-DTO mappers
- DTO-to-domain mappers
- JSON annotations

### Infrastructure Layer (3 files, ~900 lines)

**tax_jurisdiction_repository.go** (300 lines)
- PostgreSQL implementation
- Complex location queries
- Parent-child traversal
- Efficient indexing

**tax_rate_repository.go** (320 lines)
- Rate lookup by jurisdiction
- Category filtering
- Date range queries
- Bulk insert support

**tax_exemption_repository.go** (280 lines)
- Customer lookup
- Certificate validation
- Active exemption queries
- Jurisdiction filtering

### Ports Layer (1 file, ~650 lines)

**tax_handler.go** (650 lines)
- 24 REST endpoints
- Tax calculation API
- CRUD for jurisdictions
- CRUD for rates
- CRUD for exemptions
- Address validation
- Tax estimation

## Database Schema

### Tables

**blc_tax_jurisdiction**
- Hierarchical tax authorities
- Geographic location data
- Priority ordering
- Parent relationships

**blc_tax_rate**
- Tax rates with categories
- Compound tax support
- Threshold configuration
- Date range validity

**blc_tax_exemption**
- Customer exemptions
- Certificate tracking
- Jurisdiction-specific
- Category-specific

### Indexes

- Location-based queries (GIN indexes)
- Jurisdiction lookup
- Rate filtering
- Active exemptions
- Date range queries

## REST API

### Tax Calculation
```
POST /tax/calculate
POST /tax/estimate
POST /tax/validate-address
```

### Jurisdictions
```
GET    /tax/jurisdictions
POST   /tax/jurisdictions
GET    /tax/jurisdictions/{id}
PUT    /tax/jurisdictions/{id}
DELETE /tax/jurisdictions/{id}
GET    /tax/jurisdictions/code/{code}
GET    /tax/jurisdictions/country/{country}
```

### Tax Rates
```
GET    /tax/rates
POST   /tax/rates
GET    /tax/rates/{id}
PUT    /tax/rates/{id}
DELETE /tax/rates/{id}
GET    /tax/rates/jurisdiction/{id}
POST   /tax/rates/bulk
```

### Tax Exemptions
```
GET    /tax/exemptions
POST   /tax/exemptions
GET    /tax/exemptions/{id}
PUT    /tax/exemptions/{id}
DELETE /tax/exemptions/{id}
GET    /tax/exemptions/customer/{customerId}
```

## Key Features

### 1. Multi-Jurisdictional Support
- Federal, state, county, city, district levels
- Hierarchical jurisdiction relationships
- Priority-based application order
- Geographic matching

### 2. Flexible Tax Rates
- Percentage-based taxation
- Flat-rate taxation
- Compound tax calculation
- Category-specific rates
- Threshold support

### 3. Tax Exemptions
- Customer-based exemptions
- Certificate management
- Jurisdiction-specific
- Category-specific
- Date range validity

### 4. Shipping Tax
- Configurable per rate
- Category support
- Compound calculation
- Threshold support

### 5. Detailed Breakdowns
- Jurisdiction grouping
- Individual tax tracking
- Effective rate calculation
- Audit trail

### 6. Date Range Support
- Rate effective dates
- Exemption validity periods
- Historical queries
- Future planning

## Usage Examples

### Calculate Taxes for Order
```go
req := &queries.CalculateTaxRequest{
    OrderID: &orderID,
    CustomerID: &customerID,
    ShippingAddress: queries.AddressDTO{
        Country: "US",
        StateProvince: "CA",
        City: "San Francisco",
        PostalCode: "94102",
    },
    Items: []queries.TaxableItemDTO{
        {
            ItemID: "item-1",
            SKU: "WIDGET-001",
            Quantity: 2,
            UnitPrice: decimal.NewFromFloat(50.00),
            Subtotal: decimal.NewFromFloat(100.00),
            TaxCategory: "GENERAL",
        },
    },
    ShippingAmount: decimal.NewFromFloat(10.00),
}

result, err := taxCalculator.Calculate(ctx, toRequest(req))
```

### Create Tax Jurisdiction
```go
cmd := commands.CreateTaxJurisdictionCommand{
    Code: "US-CA-SF",
    Name: "San Francisco",
    JurisdictionType: "CITY",
    Country: "US",
    StateProvince: ptr("CA"),
    City: ptr("San Francisco"),
    Priority: 3,
}

jurisdiction, err := commandHandler.HandleCreateTaxJurisdiction(ctx, cmd)
```

### Create Tax Rate
```go
cmd := commands.CreateTaxRateCommand{
    JurisdictionID: jurisdictionID,
    Name: "San Francisco Sales Tax",
    TaxType: "PERCENTAGE",
    Rate: decimal.NewFromFloat(0.0625), // 6.25%
    TaxCategory: "GENERAL",
    IsShippingTaxable: true,
    Priority: 1,
}

rate, err := commandHandler.HandleCreateTaxRate(ctx, cmd)
```

### Create Tax Exemption
```go
cmd := commands.CreateTaxExemptionCommand{
    CustomerID: "customer-123",
    ExemptionCertificate: "CERT-2024-001",
    Reason: "Nonprofit organization",
    TaxCategory: ptr("GENERAL"),
}

exemption, err := commandHandler.HandleCreateTaxExemption(ctx, cmd)
```

## Testing Considerations

### Unit Tests
- Jurisdiction matching logic
- Rate calculation algorithms
- Compound tax calculation
- Exemption application
- Threshold validation

### Integration Tests
- Database queries
- Transaction handling
- Concurrent access
- Bulk operations

### API Tests
- Tax calculation endpoints
- CRUD operations
- Error handling
- Validation

## Performance Considerations

1. **Efficient Queries**
   - Indexed location fields
   - GIN indexes for arrays
   - Filtered queries
   - Batch operations

2. **Caching Opportunities**
   - Jurisdiction lookup
   - Rate configurations
   - Exemption certificates
   - Calculation results

3. **Optimization**
   - Early filtering
   - Sorted processing
   - Minimal joins
   - Bulk inserts

## Business Logic Migrated

The Tax Engine implements approximately **90%** of Broadleaf Commerce tax functionality:

### Implemented ✓
- Multi-jurisdictional tax calculation
- Tax rates with categories
- Compound tax support
- Tax exemptions
- Shipping taxation
- Threshold support
- Date range validity
- Detailed breakdowns
- REST API (24 endpoints)

### Not Implemented
- External tax provider integration (Avalara, TaxJar)
- Real-time tax rate updates
- Tax return filing
- International VAT handling
- Tax reporting

## Migration Notes

- Replaced monolithic tax module with hexagonal architecture
- Separated concerns into domain, application, infrastructure layers
- Added event-driven capabilities
- Improved testability
- Enhanced API design

## Future Enhancements

1. Tax provider integration (Avalara, TaxJar, Vertex)
2. VAT handling for international sales
3. Tax return reporting
4. Real-time rate updates
5. Tax audit trail
6. Performance optimization
7. Advanced caching strategies
8. Batch calculation APIs

## Dependencies

- `github.com/shopspring/decimal` - Precise decimal arithmetic
- `github.com/gorilla/mux` - HTTP routing
- PostgreSQL - Data persistence
- Standard library - Core functionality

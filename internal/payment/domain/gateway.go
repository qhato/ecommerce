package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// PaymentGateway defines the interface for payment gateway integrations
type PaymentGateway interface {
	// GetName returns the gateway name
	GetName() string

	// Authorize authorizes a payment (hold funds)
	Authorize(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error)

	// Capture captures a previously authorized payment
	Capture(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error)

	// Sale performs authorization and capture in one step
	Sale(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error)

	// Refund refunds a captured payment
	Refund(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error)

	// Void voids an authorized payment (before capture)
	Void(ctx context.Context, transactionID string) (*PaymentResponse, error)

	// GetTransaction retrieves transaction details
	GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error)
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID        string
	Amount         decimal.Decimal
	Currency       string
	PaymentMethod  PaymentMethod
	CardDetails    *CardDetails
	BankDetails    *BankDetails
	DigitalWallet  *DigitalWalletDetails
	BillingAddress *Address
	CustomerID     *string
	Description    string
	Metadata       map[string]string
}

// CardDetails represents credit/debit card information
type CardDetails struct {
	CardNumber     string
	ExpiryMonth    int
	ExpiryYear     int
	CVV            string
	CardholderName string
	CardType       string // VISA, MASTERCARD, AMEX, etc.
}

// BankDetails represents bank account information
type BankDetails struct {
	AccountNumber string
	RoutingNumber string
	AccountType   string // CHECKING, SAVINGS
	AccountHolder string
}

// DigitalWalletDetails represents digital wallet information
type DigitalWalletDetails struct {
	WalletType string // PAYPAL, APPLE_PAY, GOOGLE_PAY
	WalletID   string
	Email      *string
}

// Address represents a billing/shipping address
type Address struct {
	FirstName  string
	LastName   string
	Line1      string
	Line2      *string
	City       string
	State      string
	PostalCode string
	Country    string
	Phone      *string
}

// PaymentResponse represents a payment gateway response
type PaymentResponse struct {
	TransactionID   string
	GatewayResponse string
	Status          PaymentStatus
	Amount          decimal.Decimal
	Currency        string
	AuthCode        *string
	AVSResult       *string // Address Verification System result
	CVVResult       *string // Card Verification Value result
	ErrorCode       *string
	ErrorMessage    *string
	ProcessedAt     time.Time
	Metadata        map[string]string
}

// GatewayConfig represents gateway configuration
type GatewayConfig struct {
	GatewayName string
	Enabled     bool
	Priority    int    // Lower = higher priority
	Environment string // SANDBOX, PRODUCTION
	APIKey      string
	APISecret   string
	MerchantID  *string
	Config      map[string]string // Additional configuration
}

// StripeGateway implements PaymentGateway for Stripe
type StripeGateway struct {
	config *GatewayConfig
}

// NewStripeGateway creates a new Stripe gateway
func NewStripeGateway(config *GatewayConfig) *StripeGateway {
	return &StripeGateway{config: config}
}

func (g *StripeGateway) GetName() string {
	return "Stripe"
}

func (g *StripeGateway) Authorize(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	// TODO: Implement Stripe API integration
	return &PaymentResponse{
		TransactionID: "stripe_auth_123",
		Status:        PaymentStatusAuthorized,
		Amount:        request.Amount,
		Currency:      request.Currency,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *StripeGateway) Capture(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	// TODO: Implement Stripe capture
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCaptured,
		Amount:        amount,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *StripeGateway) Sale(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	// TODO: Implement Stripe sale (charge)
	return &PaymentResponse{
		TransactionID: "stripe_sale_123",
		Status:        PaymentStatusCompleted,
		Amount:        request.Amount,
		Currency:      request.Currency,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *StripeGateway) Refund(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	// TODO: Implement Stripe refund
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusRefunded,
		Amount:        amount,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *StripeGateway) Void(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// TODO: Implement Stripe void
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCancelled,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *StripeGateway) GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// TODO: Implement Stripe transaction retrieval
	return nil, nil
}

// PayPalGateway implements PaymentGateway for PayPal
type PayPalGateway struct {
	config *GatewayConfig
}

// NewPayPalGateway creates a new PayPal gateway
func NewPayPalGateway(config *GatewayConfig) *PayPalGateway {
	return &PayPalGateway{config: config}
}

func (g *PayPalGateway) GetName() string {
	return "PayPal"
}

func (g *PayPalGateway) Authorize(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	// TODO: Implement PayPal authorization
	return nil, nil
}

func (g *PayPalGateway) Capture(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	// TODO: Implement PayPal capture
	return nil, nil
}

func (g *PayPalGateway) Sale(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	// TODO: Implement PayPal sale
	return nil, nil
}

func (g *PayPalGateway) Refund(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	// TODO: Implement PayPal refund
	return nil, nil
}

func (g *PayPalGateway) Void(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// TODO: Implement PayPal void
	return nil, nil
}

func (g *PayPalGateway) GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// TODO: Implement PayPal transaction retrieval
	return nil, nil
}

// AuthorizeNetGateway implements PaymentGateway for Authorize.net
type AuthorizeNetGateway struct {
	config *GatewayConfig
}

// NewAuthorizeNetGateway creates a new Authorize.net gateway
func NewAuthorizeNetGateway(config *GatewayConfig) *AuthorizeNetGateway {
	return &AuthorizeNetGateway{config: config}
}

func (g *AuthorizeNetGateway) GetName() string {
	return "Authorize.Net"
}

func (g *AuthorizeNetGateway) Authorize(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	// TODO: Implement Authorize.net API integration
	return &PaymentResponse{
		TransactionID: "authnet_auth_123",
		Status:        PaymentStatusAuthorized,
		Amount:        request.Amount,
		Currency:      request.Currency,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *AuthorizeNetGateway) Capture(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	// TODO: Implement Authorize.net capture
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCaptured,
		Amount:        amount,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *AuthorizeNetGateway) Sale(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	// TODO: Implement Authorize.net sale
	return &PaymentResponse{
		TransactionID: "authnet_sale_123",
		Status:        PaymentStatusCompleted,
		Amount:        request.Amount,
		Currency:      request.Currency,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *AuthorizeNetGateway) Refund(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	// TODO: Implement Authorize.net refund
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusRefunded,
		Amount:        amount,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *AuthorizeNetGateway) Void(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// TODO: Implement Authorize.net void
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCancelled,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *AuthorizeNetGateway) GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// TODO: Implement Authorize.net transaction retrieval
	return nil, nil
}

// MockGateway implements PaymentGateway for testing
type MockGateway struct {
	config          *GatewayConfig
	shouldFail      bool
	failureMessage  string
}

// NewMockGateway creates a new mock gateway
func NewMockGateway(config *GatewayConfig) *MockGateway {
	return &MockGateway{config: config, shouldFail: false}
}

// SetShouldFail configures the mock to fail
func (g *MockGateway) SetShouldFail(shouldFail bool, message string) {
	g.shouldFail = shouldFail
	g.failureMessage = message
}

func (g *MockGateway) GetName() string {
	return "Mock"
}

func (g *MockGateway) Authorize(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	if g.shouldFail {
		errMsg := g.failureMessage
		return &PaymentResponse{
			TransactionID: "mock_fail_" + request.OrderID,
			Status:        PaymentStatusFailed,
			Amount:        request.Amount,
			Currency:      request.Currency,
			ErrorMessage:  &errMsg,
			ProcessedAt:   time.Now(),
		}, NewDomainError(g.failureMessage)
	}

	return &PaymentResponse{
		TransactionID: "mock_auth_" + request.OrderID,
		Status:        PaymentStatusAuthorized,
		Amount:        request.Amount,
		Currency:      request.Currency,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *MockGateway) Capture(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	if g.shouldFail {
		errMsg := g.failureMessage
		return &PaymentResponse{
			TransactionID: transactionID,
			Status:        PaymentStatusFailed,
			Amount:        amount,
			ErrorMessage:  &errMsg,
			ProcessedAt:   time.Now(),
		}, NewDomainError(g.failureMessage)
	}

	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCaptured,
		Amount:        amount,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *MockGateway) Sale(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	if g.shouldFail {
		errMsg := g.failureMessage
		return &PaymentResponse{
			TransactionID: "mock_fail_" + request.OrderID,
			Status:        PaymentStatusFailed,
			Amount:        request.Amount,
			Currency:      request.Currency,
			ErrorMessage:  &errMsg,
			ProcessedAt:   time.Now(),
		}, NewDomainError(g.failureMessage)
	}

	return &PaymentResponse{
		TransactionID: "mock_sale_" + request.OrderID,
		Status:        PaymentStatusCompleted,
		Amount:        request.Amount,
		Currency:      request.Currency,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *MockGateway) Refund(ctx context.Context, transactionID string, amount decimal.Decimal) (*PaymentResponse, error) {
	if g.shouldFail {
		errMsg := g.failureMessage
		return &PaymentResponse{
			TransactionID: transactionID,
			Status:        PaymentStatusFailed,
			Amount:        amount,
			ErrorMessage:  &errMsg,
			ProcessedAt:   time.Now(),
		}, NewDomainError(g.failureMessage)
	}

	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusRefunded,
		Amount:        amount,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *MockGateway) Void(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	if g.shouldFail {
		errMsg := g.failureMessage
		return &PaymentResponse{
			TransactionID: transactionID,
			Status:        PaymentStatusFailed,
			ErrorMessage:  &errMsg,
			ProcessedAt:   time.Now(),
		}, NewDomainError(g.failureMessage)
	}

	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCancelled,
		ProcessedAt:   time.Now(),
	}, nil
}

func (g *MockGateway) GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        PaymentStatusCompleted,
		ProcessedAt:   time.Now(),
	}, nil
}

// PaymentGatewayService manages multiple payment gateways
type PaymentGatewayService struct {
	gateways map[string]PaymentGateway
	priority []string // Gateway names in priority order
}

// NewPaymentGatewayService creates a new payment gateway service
func NewPaymentGatewayService() *PaymentGatewayService {
	return &PaymentGatewayService{
		gateways: make(map[string]PaymentGateway),
		priority: make([]string, 0),
	}
}

// RegisterGateway registers a payment gateway
func (s *PaymentGatewayService) RegisterGateway(gateway PaymentGateway) {
	s.gateways[gateway.GetName()] = gateway
	s.priority = append(s.priority, gateway.GetName())
}

// GetGateway gets a gateway by name
func (s *PaymentGatewayService) GetGateway(name string) (PaymentGateway, bool) {
	gateway, exists := s.gateways[name]
	return gateway, exists
}

// GetAllGateways returns all registered gateways
func (s *PaymentGatewayService) GetAllGateways() []PaymentGateway {
	gateways := make([]PaymentGateway, 0, len(s.priority))
	for _, name := range s.priority {
		if gateway, exists := s.gateways[name]; exists {
			gateways = append(gateways, gateway)
		}
	}
	return gateways
}

// ProcessPayment processes a payment using the specified gateway
func (s *PaymentGatewayService) ProcessPayment(
	ctx context.Context,
	gatewayName string,
	request *PaymentRequest,
) (*PaymentResponse, error) {

	gateway, exists := s.GetGateway(gatewayName)
	if !exists {
		return nil, NewDomainError("Payment gateway not found: " + gatewayName)
	}

	return gateway.Sale(ctx, request)
}

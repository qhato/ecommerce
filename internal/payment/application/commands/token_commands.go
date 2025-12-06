package commands

type CreatePaymentTokenCommand struct {
	CustomerID  string
	Token       string
	GatewayName string
	TokenType   string
	Last4Digits *string
	CardBrand   *string
	ExpiryMonth *int
	ExpiryYear  *int
	IsDefault   bool
}

type SetDefaultTokenCommand struct {
	TokenID    string
	CustomerID string
}

type DeactivateTokenCommand struct {
	TokenID string
}

type DeleteTokenCommand struct {
	TokenID string
}

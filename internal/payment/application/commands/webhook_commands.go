package commands

type ProcessWebhookCommand struct {
	GatewayName string
	EventID     string
	EventType   string
	Payload     string
	Signature   *string
	IPAddress   *string
}

type RetryWebhookCommand struct {
	WebhookID string
}

package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/qhato/ecommerce/internal/email/domain"
	"github.com/qhato/ecommerce/pkg/logger"
)

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host               string
	Port               int
	Username           string
	Password           string
	FromAddress        string
	FromName           string
	UseTLS             bool
	InsecureSkipVerify bool
	Timeout            time.Duration
}

// SMTPSender sends emails via SMTP
type SMTPSender struct {
	config *SMTPConfig
	logger logger.Logger
	queue  EmailQueue
}

// EmailQueue defines the interface for email queue
type EmailQueue interface {
	Enqueue(ctx context.Context, email *domain.Email) error
	Dequeue(ctx context.Context) (*domain.Email, error)
	Size(ctx context.Context) (int, error)
}

// NewSMTPSender creates a new SMTP sender
func NewSMTPSender(config *SMTPConfig, logger logger.Logger, queue EmailQueue) *SMTPSender {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	return &SMTPSender{
		config: config,
		logger: logger,
		queue:  queue,
	}
}

// Send sends an email immediately via SMTP
func (s *SMTPSender) Send(ctx context.Context, email *domain.Email) error {
	s.logger.Info("Sending email",
		logger.Field{Key: "email_id", Value: email.ID},
		logger.Field{Key: "to", Value: email.To},
		logger.Field{Key: "subject", Value: email.Subject},
	)

	// Mark as sending
	email.MarkAsSending()

	// Build message
	message, err := s.buildMessage(email)
	if err != nil {
		email.MarkAsFailed(fmt.Sprintf("failed to build message: %v", err))
		return fmt.Errorf("failed to build message: %w", err)
	}

	// Send via SMTP
	if err := s.sendSMTP(email.To, message); err != nil {
		email.MarkAsFailed(fmt.Sprintf("SMTP send failed: %v", err))
		s.logger.Error("Failed to send email",
			logger.Field{Key: "email_id", Value: email.ID},
			logger.Field{Key: "error", Value: err.Error()},
		)
		return fmt.Errorf("SMTP send failed: %w", err)
	}

	// Mark as sent
	email.MarkAsSent()

	s.logger.Info("Email sent successfully",
		logger.Field{Key: "email_id", Value: email.ID},
		logger.Field{Key: "to", Value: email.To},
	)

	return nil
}

// Queue queues an email for asynchronous sending
func (s *SMTPSender) Queue(ctx context.Context, email *domain.Email) error {
	s.logger.Info("Queueing email",
		logger.Field{Key: "email_id", Value: email.ID},
		logger.Field{Key: "to", Value: email.To},
	)

	if err := s.queue.Enqueue(ctx, email); err != nil {
		return fmt.Errorf("failed to enqueue email: %w", err)
	}

	return nil
}

// sendSMTP sends the email via SMTP protocol
func (s *SMTPSender) sendSMTP(to []string, message []byte) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Setup authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Send email
	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, to, message)
	}

	// Send without TLS
	return smtp.SendMail(addr, auth, s.config.FromAddress, to, message)
}

// sendWithTLS sends email with TLS
func (s *SMTPSender) sendWithTLS(addr string, auth smtp.Auth, to []string, message []byte) error {
	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName:         s.config.Host,
		InsecureSkipVerify: s.config.InsecureSkipVerify,
	}

	// Connect to SMTP server with TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with TLS: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Set sender
	if err := client.Mail(s.config.FromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data: %w", err)
	}

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	// Quit
	return client.Quit()
}

// buildMessage builds the email message
func (s *SMTPSender) buildMessage(email *domain.Email) ([]byte, error) {
	var builder strings.Builder

	// From header
	from := email.From
	if from == "" {
		from = s.config.FromAddress
	}
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, from)
	}
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from))

	// To header
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))

	// CC header
	if len(email.CC) > 0 {
		builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.CC, ", ")))
	}

	// Reply-To header
	if email.ReplyTo != "" {
		builder.WriteString(fmt.Sprintf("Reply-To: %s\r\n", email.ReplyTo))
	}

	// Subject header
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	// Custom headers
	for key, value := range email.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Date header
	builder.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))

	// MIME headers
	builder.WriteString("MIME-Version: 1.0\r\n")

	// Check if we have both plain text and HTML
	if email.Body != "" && email.HTMLBody != "" {
		boundary := fmt.Sprintf("boundary_%d", time.Now().UnixNano())
		builder.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
		builder.WriteString("\r\n")

		// Plain text part
		builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		builder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(email.Body)
		builder.WriteString("\r\n\r\n")

		// HTML part
		builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		builder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(email.HTMLBody)
		builder.WriteString("\r\n\r\n")

		builder.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else if email.HTMLBody != "" {
		// HTML only
		builder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(email.HTMLBody)
	} else {
		// Plain text only
		builder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		builder.WriteString("Content-Transfer-Encoding: 8bit\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(email.Body)
	}

	return []byte(builder.String()), nil
}

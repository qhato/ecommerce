package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"sync"
	textTemplate "text/template"
)

// TemplateRenderer renders email templates
type TemplateRenderer struct {
	templateDir   string
	htmlTemplates map[string]*template.Template
	textTemplates map[string]*textTemplate.Template
	mu            sync.RWMutex
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(templateDir string) *TemplateRenderer {
	return &TemplateRenderer{
		templateDir:   templateDir,
		htmlTemplates: make(map[string]*template.Template),
		textTemplates: make(map[string]*textTemplate.Template),
	}
}

// Render renders a template with the given data
func (r *TemplateRenderer) Render(templateName string, data map[string]interface{}) (string, string, error) {
	// Render plain text version
	plainText, err := r.renderTextTemplate(templateName, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to render text template: %w", err)
	}

	// Render HTML version
	htmlText, err := r.renderHTMLTemplate(templateName, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to render HTML template: %w", err)
	}

	return plainText, htmlText, nil
}

// renderTextTemplate renders a plain text template
func (r *TemplateRenderer) renderTextTemplate(templateName string, data map[string]interface{}) (string, error) {
	r.mu.RLock()
	tmpl, exists := r.textTemplates[templateName]
	r.mu.RUnlock()

	if !exists {
		// Load template
		var err error
		tmpl, err = r.loadTextTemplate(templateName)
		if err != nil {
			return "", err
		}

		// Cache template
		r.mu.Lock()
		r.textTemplates[templateName] = tmpl
		r.mu.Unlock()
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute text template: %w", err)
	}

	return buf.String(), nil
}

// renderHTMLTemplate renders an HTML template
func (r *TemplateRenderer) renderHTMLTemplate(templateName string, data map[string]interface{}) (string, error) {
	r.mu.RLock()
	tmpl, exists := r.htmlTemplates[templateName]
	r.mu.RUnlock()

	if !exists {
		// Load template
		var err error
		tmpl, err = r.loadHTMLTemplate(templateName)
		if err != nil {
			return "", err
		}

		// Cache template
		r.mu.Lock()
		r.htmlTemplates[templateName] = tmpl
		r.mu.Unlock()
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return buf.String(), nil
}

// loadTextTemplate loads a plain text template
func (r *TemplateRenderer) loadTextTemplate(templateName string) (*textTemplate.Template, error) {
	// Load base template
	basePath := filepath.Join(r.templateDir, "text", "base.txt")
	templatePath := filepath.Join(r.templateDir, "text", templateName+".txt")

	tmpl, err := textTemplate.New(filepath.Base(templatePath)).ParseFiles(basePath, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse text template: %w", err)
	}

	return tmpl, nil
}

// loadHTMLTemplate loads an HTML template
func (r *TemplateRenderer) loadHTMLTemplate(templateName string) (*template.Template, error) {
	// Load base template
	basePath := filepath.Join(r.templateDir, "html", "base.html")
	templatePath := filepath.Join(r.templateDir, "html", templateName+".html")

	tmpl, err := template.New(filepath.Base(templatePath)).ParseFiles(basePath, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	return tmpl, nil
}

// ClearCache clears the template cache
func (r *TemplateRenderer) ClearCache() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.htmlTemplates = make(map[string]*template.Template)
	r.textTemplates = make(map[string]*textTemplate.Template)
}

// GetDefaultFrom returns the default from address
func GetDefaultFrom(config interface{}) string {
	// This would typically come from configuration
	return "noreply@example.com"
}

// Helper functions for templates

// RenderOrderConfirmation renders order confirmation email
func RenderOrderConfirmation(data map[string]interface{}) (string, string, error) {
	renderer := NewTemplateRenderer("templates/email")
	return renderer.Render("order_confirmation", data)
}

// RenderOrderShipped renders order shipped email
func RenderOrderShipped(data map[string]interface{}) (string, string, error) {
	renderer := NewTemplateRenderer("templates/email")
	return renderer.Render("order_shipped", data)
}

// RenderPasswordReset renders password reset email
func RenderPasswordReset(data map[string]interface{}) (string, string, error) {
	renderer := NewTemplateRenderer("templates/email")
	return renderer.Render("password_reset", data)
}

// RenderWelcome renders welcome email
func RenderWelcome(data map[string]interface{}) (string, string, error) {
	renderer := NewTemplateRenderer("templates/email")
	return renderer.Render("welcome", data)
}

// RenderCartAbandonment renders cart abandonment email
func RenderCartAbandonment(data map[string]interface{}) (string, string, error) {
	renderer := NewTemplateRenderer("templates/email")
	return renderer.Render("cart_abandonment", data)
}

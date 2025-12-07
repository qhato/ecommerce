package renderer

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/qhato/ecommerce/pkg/logger"
)

// TemplateRenderer handles HTML template rendering
type TemplateRenderer struct {
	templates map[string]*template.Template
	log       *logger.Logger
	basePath  string
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer(basePath string, log *logger.Logger) (*TemplateRenderer, error) {
	r := &TemplateRenderer{
		templates: make(map[string]*template.Template),
		log:       log,
		basePath:  basePath,
	}

	// Load templates
	if err := r.loadTemplates(); err != nil {
		return nil, err
	}

	return r, nil
}

// loadTemplates loads all templates from the base path
func (r *TemplateRenderer) loadTemplates() error {
	// Admin templates
	adminLayout := filepath.Join(r.basePath, "admin/templates/layout.html")
	adminPartials := filepath.Join(r.basePath, "admin/templates/partials/*.html")

	// Dashboard
	dashboardTmpl, err := template.ParseFiles(
		adminLayout,
		filepath.Join(r.basePath, "admin/templates/dashboard.html"),
	)
	if err == nil {
		dashboardTmpl, _ = dashboardTmpl.ParseGlob(adminPartials)
		r.templates["admin/dashboard"] = dashboardTmpl
	}

	// Inventory list
	inventoryListTmpl, err := template.ParseFiles(
		adminLayout,
		filepath.Join(r.basePath, "admin/templates/inventory/list.html"),
	)
	if err == nil {
		inventoryListTmpl, _ = inventoryListTmpl.ParseGlob(adminPartials)
		r.templates["admin/inventory/list"] = inventoryListTmpl
	}

	// Storefront templates
	storefrontLayout := filepath.Join(r.basePath, "storefront/templates/layout.html")
	storefrontPartials := filepath.Join(r.basePath, "storefront/templates/partials/*.html")

	// Home
	homeTmpl, err := template.ParseFiles(
		storefrontLayout,
		filepath.Join(r.basePath, "storefront/templates/home.html"),
	)
	if err == nil {
		homeTmpl, _ = homeTmpl.ParseGlob(storefrontPartials)
		r.templates["storefront/home"] = homeTmpl
	}

	// Checkout
	checkoutTmpl, err := template.ParseFiles(
		storefrontLayout,
		filepath.Join(r.basePath, "storefront/templates/checkout.html"),
	)
	if err == nil {
		checkoutTmpl, _ = checkoutTmpl.ParseGlob(storefrontPartials)
		r.templates["storefront/checkout"] = checkoutTmpl
	}

	r.log.WithField("templates_loaded", len(r.templates)).Info("Templates loaded successfully")
	return nil
}

// Render renders a template with the given data
func (r *TemplateRenderer) Render(w io.Writer, name string, data interface{}) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}

	return tmpl.Execute(w, data)
}

// RenderHTML renders a template and writes it to the HTTP response
func (r *TemplateRenderer) RenderHTML(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.Render(w, name, data); err != nil {
		r.log.WithError(err).WithField("template", name).Error("Failed to render template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// BaseData returns base template data common to all pages
func BaseData() map[string]interface{} {
	return map[string]interface{}{
		"Year":     time.Now().Year(),
		"SiteName": "Broadleaf Commerce",
	}
}

// MergeData merges base data with page-specific data
func MergeData(base, page map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy base data
	for k, v := range base {
		result[k] = v
	}

	// Override with page data
	for k, v := range page {
		result[k] = v
	}

	return result
}

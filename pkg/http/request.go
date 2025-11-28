package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/qhato/ecommerce/pkg/errors"
)

// DecodeJSON decodes JSON request body into target
func DecodeJSON(r *http.Request, target interface{}) error {
	if r.Body == nil {
		return errors.BadRequest("Request body is empty")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return errors.BadRequest(fmt.Sprintf("Invalid JSON: %v", err))
	}

	return nil
}

// GetURLParam extracts a URL parameter from chi router
func GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// GetQueryParam extracts a query parameter
func GetQueryParam(r *http.Request, key string, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParamInt extracts an integer query parameter
func GetQueryParamInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetQueryParamBool extracts a boolean query parameter
func GetQueryParamBool(r *http.Request, key string, defaultValue bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page    int
	PerPage int
	Offset  int
}

// GetPaginationParams extracts pagination parameters from query
func GetPaginationParams(r *http.Request) PaginationParams {
	page := GetQueryParamInt(r, "page", 1)
	perPage := GetQueryParamInt(r, "per_page", 20)

	// Limit per_page to reasonable value
	if perPage > 100 {
		perPage = 100
	}
	if perPage < 1 {
		perPage = 20
	}

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * perPage

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
	}
}

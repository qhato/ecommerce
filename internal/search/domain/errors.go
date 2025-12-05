package domain

import "errors"

var (
	// ErrSearchFailed is returned when a search operation fails
	ErrSearchFailed = errors.New("search operation failed")

	// ErrIndexingFailed is returned when indexing fails
	ErrIndexingFailed = errors.New("indexing operation failed")

	// ErrDocumentNotFound is returned when a document is not found
	ErrDocumentNotFound = errors.New("document not found in index")

	// ErrInvalidQuery is returned when a search query is invalid
	ErrInvalidQuery = errors.New("invalid search query")

	// ErrSynonymNotFound is returned when a synonym is not found
	ErrSynonymNotFound = errors.New("synonym not found")

	// ErrRedirectNotFound is returned when a redirect is not found
	ErrRedirectNotFound = errors.New("search redirect not found")

	// ErrFacetConfigNotFound is returned when a facet config is not found
	ErrFacetConfigNotFound = errors.New("facet config not found")

	// ErrIndexingJobNotFound is returned when an indexing job is not found
	ErrIndexingJobNotFound = errors.New("indexing job not found")

	// ErrIndexNotAvailable is returned when the search index is not available
	ErrIndexNotAvailable = errors.New("search index not available")
)

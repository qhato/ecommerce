package persistence

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/qhato/ecommerce/internal/search/domain"
)

// PostgresSynonymRepository implements SearchSynonymRepository using PostgreSQL
type PostgresSynonymRepository struct {
	db *sql.DB
}

// NewPostgresSynonymRepository creates a new PostgreSQL synonym repository
func NewPostgresSynonymRepository(db *sql.DB) *PostgresSynonymRepository {
	return &PostgresSynonymRepository{db: db}
}

func (r *PostgresSynonymRepository) Create(synonym *domain.SearchSynonym) error {
	query := `INSERT INTO search_synonyms (term, synonyms, is_active, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, synonym.Term, pq.Array(synonym.Synonyms), synonym.IsActive,
		synonym.CreatedAt, synonym.UpdatedAt).Scan(&synonym.ID)
}

func (r *PostgresSynonymRepository) Update(synonym *domain.SearchSynonym) error {
	query := `UPDATE search_synonyms SET term=$1, synonyms=$2, is_active=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.Exec(query, synonym.Term, pq.Array(synonym.Synonyms), synonym.IsActive, synonym.UpdatedAt, synonym.ID)
	return err
}

func (r *PostgresSynonymRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM search_synonyms WHERE id=$1", id)
	return err
}

func (r *PostgresSynonymRepository) FindByID(id int64) (*domain.SearchSynonym, error) {
	var syn domain.SearchSynonym
	var synonyms pq.StringArray
	query := `SELECT id, term, synonyms, is_active, created_at, updated_at FROM search_synonyms WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&syn.ID, &syn.Term, &synonyms, &syn.IsActive, &syn.CreatedAt, &syn.UpdatedAt)
	syn.Synonyms = synonyms
	return &syn, err
}

func (r *PostgresSynonymRepository) FindByTerm(term string) (*domain.SearchSynonym, error) {
	var syn domain.SearchSynonym
	var synonyms pq.StringArray
	query := `SELECT id, term, synonyms, is_active, created_at, updated_at FROM search_synonyms WHERE term=$1`
	err := r.db.QueryRow(query, term).Scan(&syn.ID, &syn.Term, &synonyms, &syn.IsActive, &syn.CreatedAt, &syn.UpdatedAt)
	syn.Synonyms = synonyms
	return &syn, err
}

func (r *PostgresSynonymRepository) FindAll() ([]*domain.SearchSynonym, error) {
	query := `SELECT id, term, synonyms, is_active, created_at, updated_at FROM search_synonyms ORDER BY term`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var synonyms []*domain.SearchSynonym
	for rows.Next() {
		var syn domain.SearchSynonym
		var syns pq.StringArray
		if err := rows.Scan(&syn.ID, &syn.Term, &syns, &syn.IsActive, &syn.CreatedAt, &syn.UpdatedAt); err != nil {
			return nil, err
		}
		syn.Synonyms = syns
		synonyms = append(synonyms, &syn)
	}
	return synonyms, nil
}

func (r *PostgresSynonymRepository) FindActive() ([]*domain.SearchSynonym, error) {
	query := `SELECT id, term, synonyms, is_active, created_at, updated_at FROM search_synonyms WHERE is_active=true ORDER BY term`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var synonyms []*domain.SearchSynonym
	for rows.Next() {
		var syn domain.SearchSynonym
		var syns pq.StringArray
		if err := rows.Scan(&syn.ID, &syn.Term, &syns, &syn.IsActive, &syn.CreatedAt, &syn.UpdatedAt); err != nil {
			return nil, err
		}
		syn.Synonyms = syns
		synonyms = append(synonyms, &syn)
	}
	return synonyms, nil
}

// PostgresRedirectRepository implements SearchRedirectRepository using PostgreSQL
type PostgresRedirectRepository struct {
	db *sql.DB
}

// NewPostgresRedirectRepository creates a new PostgreSQL redirect repository
func NewPostgresRedirectRepository(db *sql.DB) *PostgresRedirectRepository {
	return &PostgresRedirectRepository{db: db}
}

func (r *PostgresRedirectRepository) Create(redirect *domain.SearchRedirect) error {
	query := `INSERT INTO search_redirects (search_term, target_url, priority, is_active, activation_date, expiration_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return r.db.QueryRow(query, redirect.SearchTerm, redirect.TargetURL, redirect.Priority, redirect.IsActive,
		redirect.ActivationDate, redirect.ExpirationDate, redirect.CreatedAt, redirect.UpdatedAt).Scan(&redirect.ID)
}

func (r *PostgresRedirectRepository) Update(redirect *domain.SearchRedirect) error {
	query := `UPDATE search_redirects SET search_term=$1, target_url=$2, priority=$3, is_active=$4,
              activation_date=$5, expiration_date=$6, updated_at=$7 WHERE id=$8`
	_, err := r.db.Exec(query, redirect.SearchTerm, redirect.TargetURL, redirect.Priority, redirect.IsActive,
		redirect.ActivationDate, redirect.ExpirationDate, redirect.UpdatedAt, redirect.ID)
	return err
}

func (r *PostgresRedirectRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM search_redirects WHERE id=$1", id)
	return err
}

func (r *PostgresRedirectRepository) FindByID(id int64) (*domain.SearchRedirect, error) {
	var redir domain.SearchRedirect
	query := `SELECT id, search_term, target_url, priority, is_active, activation_date, expiration_date, created_at, updated_at
              FROM search_redirects WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&redir.ID, &redir.SearchTerm, &redir.TargetURL, &redir.Priority,
		&redir.IsActive, &redir.ActivationDate, &redir.ExpirationDate, &redir.CreatedAt, &redir.UpdatedAt)
	return &redir, err
}

func (r *PostgresRedirectRepository) FindBySearchTerm(term string) (*domain.SearchRedirect, error) {
	var redir domain.SearchRedirect
	query := `SELECT id, search_term, target_url, priority, is_active, activation_date, expiration_date, created_at, updated_at
              FROM search_redirects WHERE LOWER(search_term)=LOWER($1) AND is_active=true
              ORDER BY priority DESC LIMIT 1`
	err := r.db.QueryRow(query, term).Scan(&redir.ID, &redir.SearchTerm, &redir.TargetURL, &redir.Priority,
		&redir.IsActive, &redir.ActivationDate, &redir.ExpirationDate, &redir.CreatedAt, &redir.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &redir, err
}

func (r *PostgresRedirectRepository) FindAllActive() ([]*domain.SearchRedirect, error) {
	query := `SELECT id, search_term, target_url, priority, is_active, activation_date, expiration_date, created_at, updated_at
              FROM search_redirects WHERE is_active=true ORDER BY priority DESC, search_term`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redirects []*domain.SearchRedirect
	for rows.Next() {
		var redir domain.SearchRedirect
		if err := rows.Scan(&redir.ID, &redir.SearchTerm, &redir.TargetURL, &redir.Priority, &redir.IsActive,
			&redir.ActivationDate, &redir.ExpirationDate, &redir.CreatedAt, &redir.UpdatedAt); err != nil {
			return nil, err
		}
		redirects = append(redirects, &redir)
	}
	return redirects, nil
}

// PostgresFacetConfigRepository implements SearchFacetConfigRepository using PostgreSQL
type PostgresFacetConfigRepository struct {
	db *sql.DB
}

// NewPostgresFacetConfigRepository creates a new PostgreSQL facet config repository
func NewPostgresFacetConfigRepository(db *sql.DB) *PostgresFacetConfigRepository {
	return &PostgresFacetConfigRepository{db: db}
}

func (r *PostgresFacetConfigRepository) Create(config *domain.SearchFacetConfig) error {
	query := `INSERT INTO search_facet_configs (name, label, field_name, facet_type, is_active, show_in_results,
              show_in_navigation, priority, min_doc_count, max_values, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	return r.db.QueryRow(query, config.Name, config.Label, config.FieldName, config.FacetType, config.IsActive,
		config.ShowInResults, config.ShowInNavigation, config.Priority, config.MinDocCount, config.MaxValues,
		config.CreatedAt, config.UpdatedAt).Scan(&config.ID)
}

func (r *PostgresFacetConfigRepository) Update(config *domain.SearchFacetConfig) error {
	query := `UPDATE search_facet_configs SET name=$1, label=$2, field_name=$3, facet_type=$4, is_active=$5,
              show_in_results=$6, show_in_navigation=$7, priority=$8, min_doc_count=$9, max_values=$10, updated_at=$11 WHERE id=$12`
	_, err := r.db.Exec(query, config.Name, config.Label, config.FieldName, config.FacetType, config.IsActive,
		config.ShowInResults, config.ShowInNavigation, config.Priority, config.MinDocCount, config.MaxValues,
		config.UpdatedAt, config.ID)
	return err
}

func (r *PostgresFacetConfigRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM search_facet_configs WHERE id=$1", id)
	return err
}

func (r *PostgresFacetConfigRepository) FindByID(id int64) (*domain.SearchFacetConfig, error) {
	var cfg domain.SearchFacetConfig
	query := `SELECT id, name, label, field_name, facet_type, is_active, show_in_results, show_in_navigation,
              priority, min_doc_count, max_values, created_at, updated_at FROM search_facet_configs WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&cfg.ID, &cfg.Name, &cfg.Label, &cfg.FieldName, &cfg.FacetType,
		&cfg.IsActive, &cfg.ShowInResults, &cfg.ShowInNavigation, &cfg.Priority, &cfg.MinDocCount,
		&cfg.MaxValues, &cfg.CreatedAt, &cfg.UpdatedAt)
	return &cfg, err
}

func (r *PostgresFacetConfigRepository) FindByName(name string) (*domain.SearchFacetConfig, error) {
	var cfg domain.SearchFacetConfig
	query := `SELECT id, name, label, field_name, facet_type, is_active, show_in_results, show_in_navigation,
              priority, min_doc_count, max_values, created_at, updated_at FROM search_facet_configs WHERE name=$1`
	err := r.db.QueryRow(query, name).Scan(&cfg.ID, &cfg.Name, &cfg.Label, &cfg.FieldName, &cfg.FacetType,
		&cfg.IsActive, &cfg.ShowInResults, &cfg.ShowInNavigation, &cfg.Priority, &cfg.MinDocCount,
		&cfg.MaxValues, &cfg.CreatedAt, &cfg.UpdatedAt)
	return &cfg, err
}

func (r *PostgresFacetConfigRepository) FindActive() ([]*domain.SearchFacetConfig, error) {
	query := `SELECT id, name, label, field_name, facet_type, is_active, show_in_results, show_in_navigation,
              priority, min_doc_count, max_values, created_at, updated_at FROM search_facet_configs
              WHERE is_active=true ORDER BY priority, name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*domain.SearchFacetConfig
	for rows.Next() {
		var cfg domain.SearchFacetConfig
		if err := rows.Scan(&cfg.ID, &cfg.Name, &cfg.Label, &cfg.FieldName, &cfg.FacetType, &cfg.IsActive,
			&cfg.ShowInResults, &cfg.ShowInNavigation, &cfg.Priority, &cfg.MinDocCount, &cfg.MaxValues,
			&cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, &cfg)
	}
	return configs, nil
}

func (r *PostgresFacetConfigRepository) FindForResults() ([]*domain.SearchFacetConfig, error) {
	query := `SELECT id, name, label, field_name, facet_type, is_active, show_in_results, show_in_navigation,
              priority, min_doc_count, max_values, created_at, updated_at FROM search_facet_configs
              WHERE is_active=true AND show_in_results=true ORDER BY priority, name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*domain.SearchFacetConfig
	for rows.Next() {
		var cfg domain.SearchFacetConfig
		if err := rows.Scan(&cfg.ID, &cfg.Name, &cfg.Label, &cfg.FieldName, &cfg.FacetType, &cfg.IsActive,
			&cfg.ShowInResults, &cfg.ShowInNavigation, &cfg.Priority, &cfg.MinDocCount, &cfg.MaxValues,
			&cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, &cfg)
	}
	return configs, nil
}

func (r *PostgresFacetConfigRepository) FindForNavigation() ([]*domain.SearchFacetConfig, error) {
	query := `SELECT id, name, label, field_name, facet_type, is_active, show_in_results, show_in_navigation,
              priority, min_doc_count, max_values, created_at, updated_at FROM search_facet_configs
              WHERE is_active=true AND show_in_navigation=true ORDER BY priority, name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*domain.SearchFacetConfig
	for rows.Next() {
		var cfg domain.SearchFacetConfig
		if err := rows.Scan(&cfg.ID, &cfg.Name, &cfg.Label, &cfg.FieldName, &cfg.FacetType, &cfg.IsActive,
			&cfg.ShowInResults, &cfg.ShowInNavigation, &cfg.Priority, &cfg.MinDocCount, &cfg.MaxValues,
			&cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, &cfg)
	}
	return configs, nil
}

// PostgresIndexingJobRepository implements IndexingJobRepository using PostgreSQL
type PostgresIndexingJobRepository struct {
	db *sql.DB
}

// NewPostgresIndexingJobRepository creates a new PostgreSQL indexing job repository
func NewPostgresIndexingJobRepository(db *sql.DB) *PostgresIndexingJobRepository {
	return &PostgresIndexingJobRepository{db: db}
}

func (r *PostgresIndexingJobRepository) Create(job *domain.IndexingJob) error {
	query := `INSERT INTO indexing_jobs (type, status, entity_type, total_items, processed_items, failed_items,
              error_message, started_at, completed_at, created_by, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	return r.db.QueryRow(query, job.Type, job.Status, job.EntityType, job.TotalItems, job.ProcessedItems,
		job.FailedItems, job.ErrorMessage, job.StartedAt, job.CompletedAt, job.CreatedBy, job.CreatedAt,
		job.UpdatedAt).Scan(&job.ID)
}

func (r *PostgresIndexingJobRepository) Update(job *domain.IndexingJob) error {
	query := `UPDATE indexing_jobs SET status=$1, total_items=$2, processed_items=$3, failed_items=$4,
              error_message=$5, started_at=$6, completed_at=$7, updated_at=$8 WHERE id=$9`
	_, err := r.db.Exec(query, job.Status, job.TotalItems, job.ProcessedItems, job.FailedItems,
		job.ErrorMessage, job.StartedAt, job.CompletedAt, job.UpdatedAt, job.ID)
	return err
}

func (r *PostgresIndexingJobRepository) FindByID(id int64) (*domain.IndexingJob, error) {
	var job domain.IndexingJob
	query := `SELECT id, type, status, entity_type, total_items, processed_items, failed_items, error_message,
              started_at, completed_at, created_by, created_at, updated_at FROM indexing_jobs WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&job.ID, &job.Type, &job.Status, &job.EntityType, &job.TotalItems,
		&job.ProcessedItems, &job.FailedItems, &job.ErrorMessage, &job.StartedAt, &job.CompletedAt,
		&job.CreatedBy, &job.CreatedAt, &job.UpdatedAt)
	return &job, err
}

func (r *PostgresIndexingJobRepository) FindRecent(limit int) ([]*domain.IndexingJob, error) {
	query := `SELECT id, type, status, entity_type, total_items, processed_items, failed_items, error_message,
              started_at, completed_at, created_by, created_at, updated_at FROM indexing_jobs
              ORDER BY created_at DESC LIMIT $1`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.IndexingJob
	for rows.Next() {
		var job domain.IndexingJob
		if err := rows.Scan(&job.ID, &job.Type, &job.Status, &job.EntityType, &job.TotalItems, &job.ProcessedItems,
			&job.FailedItems, &job.ErrorMessage, &job.StartedAt, &job.CompletedAt, &job.CreatedBy,
			&job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

func (r *PostgresIndexingJobRepository) FindRunning() ([]*domain.IndexingJob, error) {
	query := `SELECT id, type, status, entity_type, total_items, processed_items, failed_items, error_message,
              started_at, completed_at, created_by, created_at, updated_at FROM indexing_jobs
              WHERE status='RUNNING' ORDER BY started_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.IndexingJob
	for rows.Next() {
		var job domain.IndexingJob
		if err := rows.Scan(&job.ID, &job.Type, &job.Status, &job.EntityType, &job.TotalItems, &job.ProcessedItems,
			&job.FailedItems, &job.ErrorMessage, &job.StartedAt, &job.CompletedAt, &job.CreatedBy,
			&job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

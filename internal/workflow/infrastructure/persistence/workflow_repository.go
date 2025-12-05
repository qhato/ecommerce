package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/workflow/domain"
)

// PostgresWorkflowRepository implements domain.WorkflowRepository
type PostgresWorkflowRepository struct {
	db *sql.DB
}

// NewPostgresWorkflowRepository creates a new repository
func NewPostgresWorkflowRepository(db *sql.DB) *PostgresWorkflowRepository {
	return &PostgresWorkflowRepository{db: db}
}

// Create creates a new workflow definition
func (r *PostgresWorkflowRepository) Create(ctx context.Context, workflow *domain.Workflow) error {
	activitiesJSON, _ := json.Marshal(workflow.Activities)
	transitionsJSON, _ := json.Marshal(workflow.Transitions)
	endActivityIDsJSON, _ := json.Marshal(workflow.EndActivityIDs)
	metadataJSON, _ := json.Marshal(workflow.Metadata)

	query := `
		INSERT INTO blc_workflow (
			id, name, description, type, version, is_active,
			activities, transitions, start_activity_id, end_activity_ids,
			metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.ExecContext(ctx, query,
		workflow.ID, workflow.Name, workflow.Description, workflow.Type,
		workflow.Version, workflow.IsActive, activitiesJSON, transitionsJSON,
		workflow.StartActivityID, endActivityIDsJSON, metadataJSON,
		workflow.CreatedAt, workflow.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}

	return nil
}

// Update updates an existing workflow definition
func (r *PostgresWorkflowRepository) Update(ctx context.Context, workflow *domain.Workflow) error {
	activitiesJSON, _ := json.Marshal(workflow.Activities)
	transitionsJSON, _ := json.Marshal(workflow.Transitions)
	endActivityIDsJSON, _ := json.Marshal(workflow.EndActivityIDs)
	metadataJSON, _ := json.Marshal(workflow.Metadata)

	query := `
		UPDATE blc_workflow SET
			name = $1, description = $2, type = $3, version = $4,
			is_active = $5, activities = $6, transitions = $7,
			start_activity_id = $8, end_activity_ids = $9, metadata = $10,
			updated_at = $11
		WHERE id = $12`

	_, err := r.db.ExecContext(ctx, query,
		workflow.Name, workflow.Description, workflow.Type, workflow.Version,
		workflow.IsActive, activitiesJSON, transitionsJSON,
		workflow.StartActivityID, endActivityIDsJSON, metadataJSON,
		workflow.UpdatedAt, workflow.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
	}

	return nil
}

// FindByID finds a workflow by ID
func (r *PostgresWorkflowRepository) FindByID(ctx context.Context, id string) (*domain.Workflow, error) {
	query := `
		SELECT id, name, description, type, version, is_active,
			   activities, transitions, start_activity_id, end_activity_ids,
			   metadata, created_at, updated_at
		FROM blc_workflow WHERE id = $1`

	workflow := &domain.Workflow{}
	var activitiesJSON, transitionsJSON, endActivityIDsJSON, metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&workflow.ID, &workflow.Name, &workflow.Description, &workflow.Type,
		&workflow.Version, &workflow.IsActive, &activitiesJSON, &transitionsJSON,
		&workflow.StartActivityID, &endActivityIDsJSON, &metadataJSON,
		&workflow.CreatedAt, &workflow.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow: %w", err)
	}

	json.Unmarshal(activitiesJSON, &workflow.Activities)
	json.Unmarshal(transitionsJSON, &workflow.Transitions)
	json.Unmarshal(endActivityIDsJSON, &workflow.EndActivityIDs)
	json.Unmarshal(metadataJSON, &workflow.Metadata)

	return workflow, nil
}

// FindByName finds workflows by name
func (r *PostgresWorkflowRepository) FindByName(ctx context.Context, name string) ([]*domain.Workflow, error) {
	query := `
		SELECT id, name, description, type, version, is_active,
			   activities, transitions, start_activity_id, end_activity_ids,
			   metadata, created_at, updated_at
		FROM blc_workflow WHERE name = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflows: %w", err)
	}
	defer rows.Close()

	return r.scanWorkflows(rows)
}

// FindByType finds workflows by type
func (r *PostgresWorkflowRepository) FindByType(ctx context.Context, workflowType domain.WorkflowType, activeOnly bool) ([]*domain.Workflow, error) {
	query := `
		SELECT id, name, description, type, version, is_active,
			   activities, transitions, start_activity_id, end_activity_ids,
			   metadata, created_at, updated_at
		FROM blc_workflow WHERE type = $1`

	if activeOnly {
		query += " AND is_active = true"
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query, workflowType)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflows: %w", err)
	}
	defer rows.Close()

	return r.scanWorkflows(rows)
}

// FindAll finds all workflow definitions
func (r *PostgresWorkflowRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.Workflow, error) {
	query := `
		SELECT id, name, description, type, version, is_active,
			   activities, transitions, start_activity_id, end_activity_ids,
			   metadata, created_at, updated_at
		FROM blc_workflow`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflows: %w", err)
	}
	defer rows.Close()

	return r.scanWorkflows(rows)
}

// Delete deletes a workflow definition
func (r *PostgresWorkflowRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_workflow WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	return nil
}

// ExistsByName checks if a workflow with the given name exists
func (r *PostgresWorkflowRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM blc_workflow WHERE name = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check workflow existence: %w", err)
	}

	return exists, nil
}

// Private helper methods

func (r *PostgresWorkflowRepository) scanWorkflows(rows *sql.Rows) ([]*domain.Workflow, error) {
	workflows := make([]*domain.Workflow, 0)

	for rows.Next() {
		workflow := &domain.Workflow{}
		var activitiesJSON, transitionsJSON, endActivityIDsJSON, metadataJSON []byte

		err := rows.Scan(
			&workflow.ID, &workflow.Name, &workflow.Description, &workflow.Type,
			&workflow.Version, &workflow.IsActive, &activitiesJSON, &transitionsJSON,
			&workflow.StartActivityID, &endActivityIDsJSON, &metadataJSON,
			&workflow.CreatedAt, &workflow.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow: %w", err)
		}

		json.Unmarshal(activitiesJSON, &workflow.Activities)
		json.Unmarshal(transitionsJSON, &workflow.Transitions)
		json.Unmarshal(endActivityIDsJSON, &workflow.EndActivityIDs)
		json.Unmarshal(metadataJSON, &workflow.Metadata)

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

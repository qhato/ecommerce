package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/qhato/ecommerce/internal/workflow/domain"
)

// PostgresWorkflowExecutionRepository implements domain.WorkflowExecutionRepository
type PostgresWorkflowExecutionRepository struct {
	db *sql.DB
}

// NewPostgresWorkflowExecutionRepository creates a new repository
func NewPostgresWorkflowExecutionRepository(db *sql.DB) *PostgresWorkflowExecutionRepository {
	return &PostgresWorkflowExecutionRepository{db: db}
}

// Create creates a new workflow execution
func (r *PostgresWorkflowExecutionRepository) Create(ctx context.Context, execution *domain.WorkflowExecution) error {
	contextJSON, _ := json.Marshal(execution.Context)
	inputDataJSON, _ := json.Marshal(execution.InputData)
	outputDataJSON, _ := json.Marshal(execution.OutputData)
	activityHistoryJSON, _ := json.Marshal(execution.ActivityHistory)
	metadataJSON, _ := json.Marshal(execution.Metadata)

	query := `
		INSERT INTO blc_workflow_execution (
			id, workflow_id, workflow_version, status, context, input_data,
			output_data, current_activity_id, activity_history, error_message,
			retry_count, started_by, started_at, completed_at, last_heartbeat,
			entity_type, entity_id, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	_, err := r.db.ExecContext(ctx, query,
		execution.ID, execution.WorkflowID, execution.WorkflowVersion, execution.Status,
		contextJSON, inputDataJSON, outputDataJSON, execution.CurrentActivityID,
		activityHistoryJSON, execution.ErrorMessage, execution.RetryCount,
		execution.StartedBy, execution.StartedAt, execution.CompletedAt,
		execution.LastHeartbeat, execution.EntityType, execution.EntityID,
		metadataJSON, execution.CreatedAt, execution.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create workflow execution: %w", err)
	}

	return nil
}

// Update updates an existing workflow execution
func (r *PostgresWorkflowExecutionRepository) Update(ctx context.Context, execution *domain.WorkflowExecution) error {
	contextJSON, _ := json.Marshal(execution.Context)
	inputDataJSON, _ := json.Marshal(execution.InputData)
	outputDataJSON, _ := json.Marshal(execution.OutputData)
	activityHistoryJSON, _ := json.Marshal(execution.ActivityHistory)
	metadataJSON, _ := json.Marshal(execution.Metadata)

	query := `
		UPDATE blc_workflow_execution SET
			status = $1, context = $2, input_data = $3, output_data = $4,
			current_activity_id = $5, activity_history = $6, error_message = $7,
			retry_count = $8, completed_at = $9, last_heartbeat = $10,
			entity_type = $11, entity_id = $12, metadata = $13, updated_at = $14
		WHERE id = $15`

	_, err := r.db.ExecContext(ctx, query,
		execution.Status, contextJSON, inputDataJSON, outputDataJSON,
		execution.CurrentActivityID, activityHistoryJSON, execution.ErrorMessage,
		execution.RetryCount, execution.CompletedAt, execution.LastHeartbeat,
		execution.EntityType, execution.EntityID, metadataJSON,
		execution.UpdatedAt, execution.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update workflow execution: %w", err)
	}

	return nil
}

// FindByID finds a workflow execution by ID
func (r *PostgresWorkflowExecutionRepository) FindByID(ctx context.Context, id string) (*domain.WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, workflow_version, status, context, input_data,
			   output_data, current_activity_id, activity_history, error_message,
			   retry_count, started_by, started_at, completed_at, last_heartbeat,
			   entity_type, entity_id, metadata, created_at, updated_at
		FROM blc_workflow_execution WHERE id = $1`

	execution := &domain.WorkflowExecution{}
	var contextJSON, inputDataJSON, outputDataJSON, activityHistoryJSON, metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&execution.ID, &execution.WorkflowID, &execution.WorkflowVersion, &execution.Status,
		&contextJSON, &inputDataJSON, &outputDataJSON, &execution.CurrentActivityID,
		&activityHistoryJSON, &execution.ErrorMessage, &execution.RetryCount,
		&execution.StartedBy, &execution.StartedAt, &execution.CompletedAt,
		&execution.LastHeartbeat, &execution.EntityType, &execution.EntityID,
		&metadataJSON, &execution.CreatedAt, &execution.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow execution: %w", err)
	}

	json.Unmarshal(contextJSON, &execution.Context)
	json.Unmarshal(inputDataJSON, &execution.InputData)
	json.Unmarshal(outputDataJSON, &execution.OutputData)
	json.Unmarshal(activityHistoryJSON, &execution.ActivityHistory)
	json.Unmarshal(metadataJSON, &execution.Metadata)

	return execution, nil
}

// FindByWorkflowID finds all executions of a specific workflow
func (r *PostgresWorkflowExecutionRepository) FindByWorkflowID(ctx context.Context, workflowID string, limit int) ([]*domain.WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, workflow_version, status, context, input_data,
			   output_data, current_activity_id, activity_history, error_message,
			   retry_count, started_by, started_at, completed_at, last_heartbeat,
			   entity_type, entity_id, metadata, created_at, updated_at
		FROM blc_workflow_execution
		WHERE workflow_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, workflowID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow executions: %w", err)
	}
	defer rows.Close()

	return r.scanExecutions(rows)
}

// FindByStatus finds workflow executions by status
func (r *PostgresWorkflowExecutionRepository) FindByStatus(ctx context.Context, status domain.WorkflowStatus, limit int) ([]*domain.WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, workflow_version, status, context, input_data,
			   output_data, current_activity_id, activity_history, error_message,
			   retry_count, started_by, started_at, completed_at, last_heartbeat,
			   entity_type, entity_id, metadata, created_at, updated_at
		FROM blc_workflow_execution
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow executions: %w", err)
	}
	defer rows.Close()

	return r.scanExecutions(rows)
}

// FindByEntityReference finds workflow executions by entity reference
func (r *PostgresWorkflowExecutionRepository) FindByEntityReference(ctx context.Context, entityType, entityID string) ([]*domain.WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, workflow_version, status, context, input_data,
			   output_data, current_activity_id, activity_history, error_message,
			   retry_count, started_by, started_at, completed_at, last_heartbeat,
			   entity_type, entity_id, metadata, created_at, updated_at
		FROM blc_workflow_execution
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow executions: %w", err)
	}
	defer rows.Close()

	return r.scanExecutions(rows)
}

// FindActiveExecutions finds all active (running or suspended) workflow executions
func (r *PostgresWorkflowExecutionRepository) FindActiveExecutions(ctx context.Context, limit int) ([]*domain.WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, workflow_version, status, context, input_data,
			   output_data, current_activity_id, activity_history, error_message,
			   retry_count, started_by, started_at, completed_at, last_heartbeat,
			   entity_type, entity_id, metadata, created_at, updated_at
		FROM blc_workflow_execution
		WHERE status IN ('RUNNING', 'SUSPENDED')
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow executions: %w", err)
	}
	defer rows.Close()

	return r.scanExecutions(rows)
}

// FindStaleExecutions finds workflow executions that haven't had a heartbeat within the specified duration
func (r *PostgresWorkflowExecutionRepository) FindStaleExecutions(ctx context.Context, staleAfterMinutes int, limit int) ([]*domain.WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, workflow_version, status, context, input_data,
			   output_data, current_activity_id, activity_history, error_message,
			   retry_count, started_by, started_at, completed_at, last_heartbeat,
			   entity_type, entity_id, metadata, created_at, updated_at
		FROM blc_workflow_execution
		WHERE status = 'RUNNING'
		  AND last_heartbeat < NOW() - INTERVAL '1 minute' * $1
		ORDER BY last_heartbeat ASC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, staleAfterMinutes, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query stale executions: %w", err)
	}
	defer rows.Close()

	return r.scanExecutions(rows)
}

// Delete deletes a workflow execution
func (r *PostgresWorkflowExecutionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM blc_workflow_execution WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workflow execution: %w", err)
	}

	return nil
}

// CountByStatus counts workflow executions by status
func (r *PostgresWorkflowExecutionRepository) CountByStatus(ctx context.Context, status domain.WorkflowStatus) (int64, error) {
	query := `SELECT COUNT(*) FROM blc_workflow_execution WHERE status = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count workflow executions: %w", err)
	}

	return count, nil
}

// Private helper methods

func (r *PostgresWorkflowExecutionRepository) scanExecutions(rows *sql.Rows) ([]*domain.WorkflowExecution, error) {
	executions := make([]*domain.WorkflowExecution, 0)

	for rows.Next() {
		execution := &domain.WorkflowExecution{}
		var contextJSON, inputDataJSON, outputDataJSON, activityHistoryJSON, metadataJSON []byte

		err := rows.Scan(
			&execution.ID, &execution.WorkflowID, &execution.WorkflowVersion, &execution.Status,
			&contextJSON, &inputDataJSON, &outputDataJSON, &execution.CurrentActivityID,
			&activityHistoryJSON, &execution.ErrorMessage, &execution.RetryCount,
			&execution.StartedBy, &execution.StartedAt, &execution.CompletedAt,
			&execution.LastHeartbeat, &execution.EntityType, &execution.EntityID,
			&metadataJSON, &execution.CreatedAt, &execution.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow execution: %w", err)
		}

		json.Unmarshal(contextJSON, &execution.Context)
		json.Unmarshal(inputDataJSON, &execution.InputData)
		json.Unmarshal(outputDataJSON, &execution.OutputData)
		json.Unmarshal(activityHistoryJSON, &execution.ActivityHistory)
		json.Unmarshal(metadataJSON, &execution.Metadata)

		executions = append(executions, execution)
	}

	return executions, nil
}

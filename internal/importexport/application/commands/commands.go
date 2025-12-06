package commands

// CreateImportJobCommand creates a new import job
type CreateImportJobCommand struct {
	EntityType string `json:"entity_type"`
	Format     string `json:"format"`
	FilePath   string `json:"file_path"`
	FileName   string `json:"file_name"`
	CreatedBy  int64  `json:"created_by"`
}

// CreateExportJobCommand creates a new export job
type CreateExportJobCommand struct {
	EntityType string `json:"entity_type"`
	Format     string `json:"format"`
	CreatedBy  int64  `json:"created_by"`
}

// StartJobCommand starts a job
type StartJobCommand struct {
	ID int64 `json:"id"`
}

// CompleteJobCommand completes a job
type CompleteJobCommand struct {
	ID int64 `json:"id"`
}

// FailJobCommand marks a job as failed
type FailJobCommand struct {
	ID       int64  `json:"id"`
	ErrorLog string `json:"error_log"`
}

// UpdateProgressCommand updates job progress
type UpdateProgressCommand struct {
	ID               int64 `json:"id"`
	TotalRecords     int   `json:"total_records"`
	ProcessedRecords int   `json:"processed_records"`
	SuccessRecords   int   `json:"success_records"`
	FailedRecords    int   `json:"failed_records"`
}

// CancelJobCommand cancels a job
type CancelJobCommand struct {
	ID int64 `json:"id"`
}

// DeleteJobCommand deletes a job
type DeleteJobCommand struct {
	ID int64 `json:"id"`
}

package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/workflow/application/commands"
	"github.com/qhato/ecommerce/internal/workflow/application/queries"
	"github.com/qhato/ecommerce/internal/workflow/domain"
)

// WorkflowHandler handles HTTP requests for workflow operations
type WorkflowHandler struct {
	commandHandler *commands.WorkflowCommandHandler
	queryService   *queries.WorkflowQueryService
}

// NewWorkflowHandler creates a new workflow HTTP handler
func NewWorkflowHandler(
	commandHandler *commands.WorkflowCommandHandler,
	queryService *queries.WorkflowQueryService,
) *WorkflowHandler {
	return &WorkflowHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

// RegisterRoutes registers all workflow routes
func (h *WorkflowHandler) RegisterRoutes(router *mux.Router) {
	// Workflow Definition Endpoints
	router.HandleFunc("/workflows", h.CreateWorkflow).Methods("POST")
	router.HandleFunc("/workflows/{id}", h.GetWorkflow).Methods("GET")
	router.HandleFunc("/workflows/{id}", h.UpdateWorkflow).Methods("PUT")
	router.HandleFunc("/workflows/{id}", h.DeleteWorkflow).Methods("DELETE")
	router.HandleFunc("/workflows/{id}/activate", h.ActivateWorkflow).Methods("POST")
	router.HandleFunc("/workflows/{id}/deactivate", h.DeactivateWorkflow).Methods("POST")
	router.HandleFunc("/workflows", h.GetAllWorkflows).Methods("GET")
	router.HandleFunc("/workflows/type/{type}", h.GetWorkflowsByType).Methods("GET")
	router.HandleFunc("/workflows/name/{name}", h.GetWorkflowsByName).Methods("GET")

	// Workflow Execution Endpoints
	router.HandleFunc("/workflows/{id}/execute", h.StartWorkflowExecution).Methods("POST")
	router.HandleFunc("/workflow-executions/{id}", h.GetWorkflowExecution).Methods("GET")
	router.HandleFunc("/workflow-executions", h.GetWorkflowExecutions).Methods("GET")
	router.HandleFunc("/workflow-executions/{id}/suspend", h.SuspendWorkflow).Methods("POST")
	router.HandleFunc("/workflow-executions/{id}/resume", h.ResumeWorkflow).Methods("POST")
	router.HandleFunc("/workflow-executions/{id}/cancel", h.CancelWorkflow).Methods("POST")
	router.HandleFunc("/workflow-executions/{id}/context", h.SetWorkflowContext).Methods("POST")

	// Activity Execution Endpoints
	router.HandleFunc("/workflow-executions/{id}/activities/{activityId}/complete", h.CompleteActivity).Methods("POST")
	router.HandleFunc("/workflow-executions/{id}/activities/{activityId}/fail", h.FailActivity).Methods("POST")

	// Query Endpoints
	router.HandleFunc("/workflow-executions/by-entity", h.GetExecutionsByEntity).Methods("GET")
	router.HandleFunc("/workflow-executions/active", h.GetActiveExecutions).Methods("GET")
	router.HandleFunc("/workflow-executions/stale", h.GetStaleExecutions).Methods("GET")
}

// CreateWorkflow handles creating a new workflow definition
func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	workflow, err := h.commandHandler.HandleCreateWorkflow(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowDTO(workflow)
	respondWithJSON(w, http.StatusCreated, response)
}

// GetWorkflow handles retrieving a workflow by ID
func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	workflow, err := h.queryService.GetWorkflow(r.Context(), id)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, workflow)
}

// UpdateWorkflow handles updating a workflow definition
func (h *WorkflowHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var cmd commands.UpdateWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.ID = id

	workflow, err := h.commandHandler.HandleUpdateWorkflow(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowDTO(workflow)
	respondWithJSON(w, http.StatusOK, response)
}

// DeleteWorkflow handles deleting a workflow definition
func (h *WorkflowHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.DeleteWorkflowCommand{ID: id}

	if err := h.commandHandler.HandleDeleteWorkflow(r.Context(), cmd); err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// ActivateWorkflow handles activating a workflow
func (h *WorkflowHandler) ActivateWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.ActivateWorkflowCommand{ID: id}

	workflow, err := h.commandHandler.HandleActivateWorkflow(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowDTO(workflow)
	respondWithJSON(w, http.StatusOK, response)
}

// DeactivateWorkflow handles deactivating a workflow
func (h *WorkflowHandler) DeactivateWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmd := commands.DeactivateWorkflowCommand{ID: id}

	workflow, err := h.commandHandler.HandleDeactivateWorkflow(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowDTO(workflow)
	respondWithJSON(w, http.StatusOK, response)
}

// GetAllWorkflows handles retrieving all workflows
func (h *WorkflowHandler) GetAllWorkflows(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"

	workflows, err := h.queryService.GetAllWorkflows(r.Context(), activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get workflows: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workflows)
}

// GetWorkflowsByType handles retrieving workflows by type
func (h *WorkflowHandler) GetWorkflowsByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workflowType := vars["type"]
	activeOnly := r.URL.Query().Get("active") == "true"

	workflows, err := h.queryService.GetWorkflowsByType(r.Context(), workflowType, activeOnly)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get workflows: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workflows)
}

// GetWorkflowsByName handles retrieving workflows by name
func (h *WorkflowHandler) GetWorkflowsByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	workflows, err := h.queryService.GetWorkflowsByName(r.Context(), name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get workflows: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, workflows)
}

// StartWorkflowExecution handles starting a new workflow execution
func (h *WorkflowHandler) StartWorkflowExecution(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workflowID := vars["id"]

	var cmd commands.StartWorkflowExecutionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cmd.WorkflowID = workflowID

	execution, err := h.commandHandler.HandleStartWorkflowExecution(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowExecutionDTO(execution)
	respondWithJSON(w, http.StatusCreated, response)
}

// GetWorkflowExecution handles retrieving a workflow execution
func (h *WorkflowHandler) GetWorkflowExecution(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	execution, err := h.queryService.GetWorkflowExecution(r.Context(), id)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, execution)
}

// GetWorkflowExecutions handles retrieving workflow executions with filters
func (h *WorkflowHandler) GetWorkflowExecutions(w http.ResponseWriter, r *http.Request) {
	workflowID := r.URL.Query().Get("workflow_id")
	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	var executions []*queries.WorkflowExecutionDTO
	var err error

	if workflowID != "" {
		executions, err = h.queryService.GetWorkflowExecutionsByWorkflowID(r.Context(), workflowID, limit)
	} else if status != "" {
		executions, err = h.queryService.GetWorkflowExecutionsByStatus(r.Context(), status, limit)
	} else {
		executions, err = h.queryService.GetActiveWorkflowExecutions(r.Context(), limit)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get executions: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, executions)
}

// GetExecutionsByEntity handles retrieving executions by entity reference
func (h *WorkflowHandler) GetExecutionsByEntity(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("entity_type")
	entityID := r.URL.Query().Get("entity_id")

	if entityType == "" || entityID == "" {
		respondWithError(w, http.StatusBadRequest, "entity_type and entity_id are required")
		return
	}

	executions, err := h.queryService.GetWorkflowExecutionsByEntityReference(r.Context(), entityType, entityID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get executions: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, executions)
}

// GetActiveExecutions handles retrieving active workflow executions
func (h *WorkflowHandler) GetActiveExecutions(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	executions, err := h.queryService.GetActiveWorkflowExecutions(r.Context(), limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get active executions: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, executions)
}

// GetStaleExecutions handles retrieving stale workflow executions
func (h *WorkflowHandler) GetStaleExecutions(w http.ResponseWriter, r *http.Request) {
	staleAfterStr := r.URL.Query().Get("stale_after_minutes")
	limitStr := r.URL.Query().Get("limit")

	staleAfter := 30
	if staleAfterStr != "" {
		if s, err := strconv.Atoi(staleAfterStr); err == nil {
			staleAfter = s
		}
	}

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	executions, err := h.queryService.GetStaleWorkflowExecutions(r.Context(), staleAfter, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get stale executions: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, executions)
}

// SuspendWorkflow handles suspending a workflow execution
func (h *WorkflowHandler) SuspendWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		SuspendReason string `json:"suspend_reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.SuspendWorkflowCommand{
		ExecutionID:   id,
		SuspendReason: req.SuspendReason,
	}

	execution, err := h.commandHandler.HandleSuspendWorkflow(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowExecutionDTO(execution)
	respondWithJSON(w, http.StatusOK, response)
}

// ResumeWorkflow handles resuming a suspended workflow execution
func (h *WorkflowHandler) ResumeWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		ResumedBy string `json:"resumed_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.ResumeWorkflowCommand{
		ExecutionID: id,
		ResumedBy:   req.ResumedBy,
	}

	execution, err := h.commandHandler.HandleResumeWorkflow(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowExecutionDTO(execution)
	respondWithJSON(w, http.StatusOK, response)
}

// CancelWorkflow handles cancelling a workflow execution
func (h *WorkflowHandler) CancelWorkflow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		CancelReason string `json:"cancel_reason"`
		CancelledBy  string `json:"cancelled_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.CancelWorkflowCommand{
		ExecutionID:  id,
		CancelReason: req.CancelReason,
		CancelledBy:  req.CancelledBy,
	}

	execution, err := h.commandHandler.HandleCancelWorkflow(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowExecutionDTO(execution)
	respondWithJSON(w, http.StatusOK, response)
}

// SetWorkflowContext handles setting a value in workflow context
func (h *WorkflowHandler) SetWorkflowContext(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.SetWorkflowContextCommand{
		ExecutionID: id,
		Key:         req.Key,
		Value:       req.Value,
	}

	execution, err := h.commandHandler.HandleSetWorkflowContext(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowExecutionDTO(execution)
	respondWithJSON(w, http.StatusOK, response)
}

// CompleteActivity handles completing an activity
func (h *WorkflowHandler) CompleteActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	executionID := vars["id"]
	activityID := vars["activityId"]

	var req struct {
		OutputData map[string]interface{} `json:"output_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.CompleteActivityCommand{
		ExecutionID: executionID,
		ActivityID:  activityID,
		OutputData:  req.OutputData,
	}

	execution, err := h.commandHandler.HandleCompleteActivity(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowExecutionDTO(execution)
	respondWithJSON(w, http.StatusOK, response)
}

// FailActivity handles failing an activity
func (h *WorkflowHandler) FailActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	executionID := vars["id"]
	activityID := vars["activityId"]

	var req struct {
		ErrorMessage string `json:"error_message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cmd := commands.FailActivityCommand{
		ExecutionID:  executionID,
		ActivityID:   activityID,
		ErrorMessage: req.ErrorMessage,
	}

	execution, err := h.commandHandler.HandleFailActivity(r.Context(), cmd)
	if err != nil {
		h.handleWorkflowError(w, err)
		return
	}

	response := queries.ToWorkflowExecutionDTO(execution)
	respondWithJSON(w, http.StatusOK, response)
}

// Helper methods

func (h *WorkflowHandler) handleWorkflowError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrWorkflowNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case domain.ErrWorkflowExecutionNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	case domain.ErrWorkflowAlreadyExists:
		respondWithError(w, http.StatusConflict, err.Error())
	case domain.ErrWorkflowInactive:
		respondWithError(w, http.StatusPreconditionFailed, err.Error())
	case domain.ErrWorkflowNotRunning:
		respondWithError(w, http.StatusPreconditionFailed, err.Error())
	case domain.ErrWorkflowNotSuspended:
		respondWithError(w, http.StatusPreconditionFailed, err.Error())
	case domain.ErrWorkflowAlreadyCompleted:
		respondWithError(w, http.StatusConflict, err.Error())
	case domain.ErrWorkflowAlreadyCancelled:
		respondWithError(w, http.StatusConflict, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Workflow operation failed: "+err.Error())
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

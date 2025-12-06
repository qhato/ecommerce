package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/importexport/application/commands"
	"github.com/qhato/ecommerce/internal/importexport/application/queries"
	"github.com/qhato/ecommerce/internal/importexport/domain"
)

type JobHandler struct {
	commandHandler *commands.JobCommandHandler
	queryService   *queries.JobQueryService
}

func NewJobHandler(
	commandHandler *commands.JobCommandHandler,
	queryService *queries.JobQueryService,
) *JobHandler {
	return &JobHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *JobHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import-export/jobs", h.CreateImportJob).Methods("POST").Queries("type", "import")
	router.HandleFunc("/import-export/jobs", h.CreateExportJob).Methods("POST").Queries("type", "export")
	router.HandleFunc("/import-export/jobs/{id}", h.GetJob).Methods("GET")
	router.HandleFunc("/import-export/jobs", h.GetRecentJobs).Methods("GET")
	router.HandleFunc("/import-export/jobs/type/{type}", h.GetJobsByType).Methods("GET")
	router.HandleFunc("/import-export/jobs/status/{status}", h.GetJobsByStatus).Methods("GET")
	router.HandleFunc("/import-export/jobs/{id}/start", h.StartJob).Methods("POST")
	router.HandleFunc("/import-export/jobs/{id}/complete", h.CompleteJob).Methods("POST")
	router.HandleFunc("/import-export/jobs/{id}/fail", h.FailJob).Methods("POST")
	router.HandleFunc("/import-export/jobs/{id}/progress", h.UpdateProgress).Methods("PUT")
	router.HandleFunc("/import-export/jobs/{id}/cancel", h.CancelJob).Methods("POST")
	router.HandleFunc("/import-export/jobs/{id}", h.DeleteJob).Methods("DELETE")
}

func (h *JobHandler) CreateImportJob(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateImportJobCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job, err := h.commandHandler.HandleCreateImportJob(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(queries.ToJobDTO(job))
}

func (h *JobHandler) CreateExportJob(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateExportJobCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job, err := h.commandHandler.HandleCreateExportJob(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(queries.ToJobDTO(job))
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := h.queryService.GetJob(r.Context(), id)
	if err != nil {
		if err == domain.ErrJobNotFound {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) GetRecentJobs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	jobs, err := h.queryService.GetRecentJobs(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) GetJobsByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobType := vars["type"]

	status := r.URL.Query().Get("status")
	if status == "" {
		status = string(domain.JobStatusPending)
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	jobs, err := h.queryService.GetJobsByType(r.Context(), jobType, status, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) GetJobsByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	jobs, err := h.queryService.GetJobsByStatus(r.Context(), status, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) StartJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	cmd := commands.StartJobCommand{ID: id}
	job, err := h.commandHandler.HandleStartJob(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJobNotFound {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		if err == domain.ErrJobAlreadyRunning {
			http.Error(w, "Job is already running", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToJobDTO(job))
}

func (h *JobHandler) CompleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	cmd := commands.CompleteJobCommand{ID: id}
	job, err := h.commandHandler.HandleCompleteJob(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJobNotFound {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		if err == domain.ErrJobNotProcessing {
			http.Error(w, "Job is not in processing state", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToJobDTO(job))
}

func (h *JobHandler) FailJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var cmd commands.FailJobCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	job, err := h.commandHandler.HandleFailJob(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJobNotFound {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToJobDTO(job))
}

func (h *JobHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateProgressCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	job, err := h.commandHandler.HandleUpdateProgress(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJobNotFound {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToJobDTO(job))
}

func (h *JobHandler) CancelJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	cmd := commands.CancelJobCommand{ID: id}
	job, err := h.commandHandler.HandleCancelJob(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrJobNotFound {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queries.ToJobDTO(job))
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeleteJobCommand{ID: id}
	if err := h.commandHandler.HandleDeleteJob(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/return/application/commands"
	"github.com/qhato/ecommerce/internal/return/application/queries"
	"github.com/qhato/ecommerce/internal/return/domain"
)

type ReturnHandler struct {
	commandHandler *commands.ReturnCommandHandler
	queryService   *queries.ReturnQueryService
}

func NewReturnHandler(
	commandHandler *commands.ReturnCommandHandler,
	queryService *queries.ReturnQueryService,
) *ReturnHandler {
	return &ReturnHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
	}
}

func (h *ReturnHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/returns", h.CreateReturn).Methods("POST")
	router.HandleFunc("/returns/{id}", h.GetReturn).Methods("GET")
	router.HandleFunc("/returns/rma/{rma}", h.GetReturnByRMA).Methods("GET")
	router.HandleFunc("/returns/customer/{customerId}", h.GetReturnsByCustomer).Methods("GET")
	router.HandleFunc("/returns/status/{status}", h.GetReturnsByStatus).Methods("GET")
	router.HandleFunc("/returns/{id}/approve", h.ApproveReturn).Methods("POST")
	router.HandleFunc("/returns/{id}/reject", h.RejectReturn).Methods("POST")
	router.HandleFunc("/returns/{id}/receive", h.ReceiveReturn).Methods("POST")
	router.HandleFunc("/returns/{id}/inspect", h.InspectReturn).Methods("POST")
	router.HandleFunc("/returns/{id}/refund", h.ProcessRefund).Methods("POST")
	router.HandleFunc("/returns/{id}/cancel", h.CancelReturn).Methods("POST")
}

func (h *ReturnHandler) CreateReturn(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateReturnCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	returnReq, err := h.commandHandler.HandleCreateReturn(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) GetReturn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid return ID", http.StatusBadRequest)
		return
	}

	returnReq, err := h.queryService.GetReturn(r.Context(), id)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) GetReturnByRMA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rma := vars["rma"]

	returnReq, err := h.queryService.GetReturnByRMA(r.Context(), rma)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) GetReturnsByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customerId"]

	returns, err := h.queryService.GetReturnsByCustomer(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returns)
}

func (h *ReturnHandler) GetReturnsByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	returns, err := h.queryService.GetReturnsByStatus(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returns)
}

func (h *ReturnHandler) ApproveReturn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid return ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ApproveReturnCommand{ID: id}
	returnReq, err := h.commandHandler.HandleApproveReturn(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) RejectReturn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid return ID", http.StatusBadRequest)
		return
	}

	var cmd commands.RejectReturnCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	returnReq, err := h.commandHandler.HandleRejectReturn(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) ReceiveReturn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid return ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ReceiveReturnCommand{ID: id}
	returnReq, err := h.commandHandler.HandleReceiveReturn(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) InspectReturn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid return ID", http.StatusBadRequest)
		return
	}

	var cmd commands.InspectReturnCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	returnReq, err := h.commandHandler.HandleInspectReturn(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid return ID", http.StatusBadRequest)
		return
	}

	var cmd commands.ProcessRefundCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	returnReq, err := h.commandHandler.HandleProcessRefund(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

func (h *ReturnHandler) CancelReturn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid return ID", http.StatusBadRequest)
		return
	}

	cmd := commands.CancelReturnCommand{ID: id}
	returnReq, err := h.commandHandler.HandleCancelReturn(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrReturnNotFound {
			http.Error(w, "Return not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnReq)
}

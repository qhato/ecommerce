package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/qhato/ecommerce/internal/admin/application/commands"
	"github.com/qhato/ecommerce/internal/admin/application/queries"
	"github.com/qhato/ecommerce/internal/admin/application/services"
	"github.com/qhato/ecommerce/internal/admin/domain"
)

type AdminHandler struct {
	commandHandler *commands.AdminCommandHandler
	queryService   *queries.AdminQueryService
	authService    *services.AuthenticationService
	authzService   *services.AuthorizationService
}

func NewAdminHandler(
	commandHandler *commands.AdminCommandHandler,
	queryService *queries.AdminQueryService,
	authService *services.AuthenticationService,
	authzService *services.AuthorizationService,
) *AdminHandler {
	return &AdminHandler{
		commandHandler: commandHandler,
		queryService:   queryService,
		authService:    authService,
		authzService:   authzService,
	}
}

func (h *AdminHandler) RegisterRoutes(router *mux.Router) {
	// Authentication endpoints
	router.HandleFunc("/admin/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/admin/auth/logout", h.Logout).Methods("POST")
	router.HandleFunc("/admin/auth/refresh", h.RefreshToken).Methods("POST")
	router.HandleFunc("/admin/auth/me", h.GetCurrentUser).Methods("GET")

	// User management endpoints
	router.HandleFunc("/admin/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/admin/users", h.GetAllUsers).Methods("GET")
	router.HandleFunc("/admin/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/admin/users/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/admin/users/{id}/activate", h.ActivateUser).Methods("POST")
	router.HandleFunc("/admin/users/{id}/deactivate", h.DeactivateUser).Methods("POST")
	router.HandleFunc("/admin/users/{id}/password", h.ChangePassword).Methods("PUT")
	router.HandleFunc("/admin/users/{id}/roles", h.AssignRole).Methods("POST")
	router.HandleFunc("/admin/users/{id}/roles/{roleId}", h.UnassignRole).Methods("DELETE")

	// Role management endpoints
	router.HandleFunc("/admin/roles", h.CreateRole).Methods("POST")
	router.HandleFunc("/admin/roles", h.GetAllRoles).Methods("GET")
	router.HandleFunc("/admin/roles/{id}", h.GetRole).Methods("GET")
	router.HandleFunc("/admin/roles/{id}", h.UpdateRole).Methods("PUT")
	router.HandleFunc("/admin/roles/{id}/permissions", h.GrantPermission).Methods("POST")
	router.HandleFunc("/admin/roles/{id}/permissions/{permissionId}", h.RevokePermission).Methods("DELETE")

	// Permission management endpoints
	router.HandleFunc("/admin/permissions", h.CreatePermission).Methods("POST")
	router.HandleFunc("/admin/permissions", h.GetAllPermissions).Methods("GET")
	router.HandleFunc("/admin/permissions/{id}", h.GetPermission).Methods("GET")
	router.HandleFunc("/admin/permissions/{id}", h.UpdatePermission).Methods("PUT")
	router.HandleFunc("/admin/permissions/resource/{resource}", h.GetPermissionsByResource).Methods("GET")

	// Audit log endpoints
	router.HandleFunc("/admin/audit-logs", h.GetRecentAuditLogs).Methods("GET")
	router.HandleFunc("/admin/audit-logs/{id}", h.GetAuditLog).Methods("GET")
	router.HandleFunc("/admin/audit-logs/user/{userId}", h.GetAuditLogsByUser).Methods("GET")
	router.HandleFunc("/admin/audit-logs/security-events", h.GetSecurityEvents).Methods("GET")
	router.HandleFunc("/admin/audit-logs/failed-logins/{username}", h.GetFailedLogins).Methods("GET")

	// Authorization check endpoint
	router.HandleFunc("/admin/authorize", h.CheckPermission).Methods("POST")
}

// Authentication handlers

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ipAddress := h.getIPAddress(r)
	userAgent := r.UserAgent()

	tokens, user, err := h.authService.Login(r.Context(), req.Username, req.Password, ipAddress, userAgent)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_in":    tokens.ExpiresIn,
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"is_super":   user.IsSuper,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT tokens are stateless, so logout is handled client-side
	// Client should delete the token from storage
	// In a production system, you might want to maintain a blacklist of tokens
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *AdminHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	token := h.extractToken(r)
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	claims, err := h.authService.ValidateToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Load full user details
	userDTO, err := h.queryService.GetUser(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDTO)
}

// User management handlers

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateUserCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.commandHandler.HandleCreateUser(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrUsernameTaken || err == domain.ErrEmailTaken {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	users, err := h.queryService.GetAllUsers(r.Context(), activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.queryService.GetUser(r.Context(), id)
	if err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateUserCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	user, err := h.commandHandler.HandleUpdateUser(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AdminHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cmd := commands.ActivateUserCommand{UserID: id}
	user, err := h.commandHandler.HandleActivateUser(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AdminHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cmd := commands.DeactivateUserCommand{UserID: id}
	user, err := h.commandHandler.HandleDeactivateUser(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AdminHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var cmd commands.ChangePasswordCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.UserID = id

	if err := h.commandHandler.HandleChangePassword(r.Context(), cmd); err != nil {
		if err == domain.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		RoleID int64 `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cmd := commands.AssignRoleCommand{
		UserID: userID,
		RoleID: req.RoleID,
	}

	if err := h.commandHandler.HandleAssignRole(r.Context(), cmd); err != nil {
		if err == domain.ErrUserNotFound || err == domain.ErrRoleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) UnassignRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(vars["roleId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	cmd := commands.UnassignRoleCommand{
		UserID: userID,
		RoleID: roleID,
	}

	if err := h.commandHandler.HandleUnassignRole(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Role management handlers

func (h *AdminHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateRoleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role, err := h.commandHandler.HandleCreateRole(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrRoleNameTaken {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func (h *AdminHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	roles, err := h.queryService.GetAllRoles(r.Context(), activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func (h *AdminHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	role, err := h.queryService.GetRole(r.Context(), id)
	if err != nil {
		if err == domain.ErrRoleNotFound {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *AdminHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdateRoleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	role, err := h.commandHandler.HandleUpdateRole(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrRoleNotFound {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *AdminHandler) GrantPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var req struct {
		PermissionID int64 `json:"permission_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cmd := commands.GrantPermissionCommand{
		RoleID:       roleID,
		PermissionID: req.PermissionID,
	}

	if err := h.commandHandler.HandleGrantPermission(r.Context(), cmd); err != nil {
		if err == domain.ErrRoleNotFound || err == domain.ErrPermissionNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) RevokePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	permissionID, err := strconv.ParseInt(vars["permissionId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	cmd := commands.RevokePermissionCommand{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	if err := h.commandHandler.HandleRevokePermission(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Permission management handlers

func (h *AdminHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreatePermissionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	permission, err := h.commandHandler.HandleCreatePermission(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrPermissionAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(permission)
}

func (h *AdminHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	permissions, err := h.queryService.GetAllPermissions(r.Context(), activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func (h *AdminHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.queryService.GetPermission(r.Context(), id)
	if err != nil {
		if err == domain.ErrPermissionNotFound {
			http.Error(w, "Permission not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

func (h *AdminHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	var cmd commands.UpdatePermissionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	permission, err := h.commandHandler.HandleUpdatePermission(r.Context(), cmd)
	if err != nil {
		if err == domain.ErrPermissionNotFound {
			http.Error(w, "Permission not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permission)
}

func (h *AdminHandler) GetPermissionsByResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]
	activeOnly := r.URL.Query().Get("active_only") == "true"

	permissions, err := h.queryService.GetPermissionsByResource(r.Context(), resource, activeOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

// Audit log handlers

func (h *AdminHandler) GetRecentAuditLogs(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.queryService.GetRecentAuditLogs(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *AdminHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid audit log ID", http.StatusBadRequest)
		return
	}

	log, err := h.queryService.GetAuditLog(r.Context(), id)
	if err != nil {
		if err == domain.ErrAuditLogNotFound {
			http.Error(w, "Audit log not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(log)
}

func (h *AdminHandler) GetAuditLogsByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.queryService.GetAuditLogsByUser(r.Context(), userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *AdminHandler) GetSecurityEvents(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.queryService.GetSecurityEvents(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *AdminHandler) GetFailedLogins(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.queryService.GetFailedLogins(r.Context(), username, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// Authorization handler

func (h *AdminHandler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int64  `json:"user_id"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hasPermission, err := h.authzService.CheckResourcePermission(
		r.Context(),
		req.UserID,
		domain.PermissionResource(req.Resource),
		domain.PermissionAction(req.Action),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"has_permission": hasPermission,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper methods

func (h *AdminHandler) extractToken(r *http.Request) string {
	// Try to get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Try to get token from query parameter
	return r.URL.Query().Get("token")
}

func (h *AdminHandler) getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Use RemoteAddr
	return r.RemoteAddr
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sandeepshukla/server2admin/client"
	"github.com/sandeepshukla/server2admin/models"
)

type AdminHandler struct {
	userClient *client.UserClient
}

func NewAdminHandler(c *client.UserClient) *AdminHandler {
	return &AdminHandler{userClient: c}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}

// GET /admin/users
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userClient.GetAllUsers()
	if err != nil {
		writeError(w, http.StatusBadGateway, "could not fetch users from Server1User: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// GET /admin/users/{id}
func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/admin/users/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	user, status, err := h.userClient.GetUserByID(id)
	if err != nil {
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// POST /admin/users
func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Age <= 0 {
		writeError(w, http.StatusBadRequest, "age must be a positive number")
		return
	}

	user, status, err := h.userClient.CreateUser(req)
	if err != nil {
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

// PUT /admin/users/{id}
func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/admin/users/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Age <= 0 {
		writeError(w, http.StatusBadRequest, "age must be a positive number")
		return
	}

	user, status, err := h.userClient.UpdateUser(id, req)
	if err != nil {
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// DELETE /admin/users/{id}
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/admin/users/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing user id")
		return
	}

	status, err := h.userClient.DeleteUser(id)
	if err != nil {
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "user deleted via Server1User"})
}

func extractID(path, prefix string) string {
	id := strings.TrimPrefix(path, prefix)
	return strings.TrimSuffix(id, "/")
}

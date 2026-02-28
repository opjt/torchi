package handler

import (
	"encoding/json"
	"net/http"
	"torchi/internal/infrastructure/db/postgresql"
)

type HealthHandler struct {
	db *postgresql.Database
}

func NewHealthHandler(db *postgresql.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "ok"
	httpStatus := http.StatusOK

	if err := h.db.Ping(r.Context()); err != nil {
		dbStatus = "unreachable"
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]string{
		"status": status,
		"db":     dbStatus,
	})
}

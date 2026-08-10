package handlers

import (
	"net/http"

	m "github.com/Diaku49/AI-visibility-tracker/internal/api/middleware"
	"github.com/google/uuid"
)

func (h *ServerHandler) BeginScan(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		if writeErr := HTTPError(w, http.StatusForbidden, err.Error(), nil); writeErr != nil {
			l.Error("write begin scan error response", "error", writeErr)
		}
		l.Error("get user ID for begin scan", "error", err)
		return
	}

	projectIDStr := r.PathValue("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusBadRequest, "invalid project ID", nil); writeErr != nil {
			l.Error("write begin scan error response", "error", writeErr)
		}
		l.Error("parse project ID for begin scan", "error", err)
		return
	}

	scanID, err := h.st.CreateScan(r.Context(), userID, projectID)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, "failed to create scan", nil); writeErr != nil {
			l.Error("write begin scan error response", "error", writeErr)
		}
		l.Error("create scan", "error", err)
		return
	}

	if err := HTTPResponse(w, http.StatusCreated, "scan created and its pending", nil); err != nil {
		l.Error("write begin scan response", "error", err)
	}
	l.Info("scan created and its pending", "project_id", projectID, "scan_id", scanID)
}

func (h *ServerHandler) GetAllScansForUser(w http.ResponseWriter, r *http.Request) {
}

func (h *ServerHandler) GetScanByID(w http.ResponseWriter, r *http.Request) {}

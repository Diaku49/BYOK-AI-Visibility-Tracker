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
		HTTPError(w, http.StatusForbidden, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	projectIDStr := r.PathValue("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "invalid project ID", nil)
		l.Error(err.Error())
		return
	}

	scanID, err := h.st.CreateScan(r.Context(), userID, projectID)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, "failed to create scan", nil)
		l.Error(err.Error())
		return
	}

	HTTPResponse(w, http.StatusCreated, "scan created and its pending", nil)
	l.Info("scan created and its pending", "project_id", projectID, "scan_id", scanID)
}

func (h *ServerHandler) GetAllScansForUser(w http.ResponseWriter, r *http.Request) {
}

func (h *ServerHandler) GetScanByID(w http.ResponseWriter, r *http.Request) {}

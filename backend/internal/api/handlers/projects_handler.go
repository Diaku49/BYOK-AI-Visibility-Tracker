package handlers

import (
	"net/http"

	m "github.com/Diaku49/AI-visibility-tracker/internal/api/middleware"
	"github.com/Diaku49/AI-visibility-tracker/internal/dto"
	"github.com/google/uuid"
)

func (h *ServerHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	var req dto.CreateProject

	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		if writeErr := HTTPError(w, http.StatusForbidden, err.Error(), nil); writeErr != nil {
			l.Error("write create project error response", "error", writeErr)
		}
		l.Error("get user ID for create project", "error", err)
		return
	}

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		if writeErr := HTTPError(w, http.StatusBadRequest, err.Error(), nil); writeErr != nil {
			l.Error("write create project error response", "error", writeErr)
		}
		l.Error("create project request validation failed", "error", err)
		return
	}

	projectID, err := h.st.CreateProject(r.Context(), userID, req)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, "failed to create project", nil); writeErr != nil {
			l.Error("write create project error response", "error", writeErr)
		}
		l.Error("create project", "error", err)
		return
	}

	if err := HTTPResponse(w, http.StatusCreated, "project created successfully", map[string]string{"project_id": projectID.String()}); err != nil {
		l.Error("write create project response", "error", err)
	}
	l.Info("project created successfully", "user_id", userID.String(), "project_id", projectID.String())
}

func (h *ServerHandler) GetAllProjectsForUser(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		if writeErr := HTTPError(w, http.StatusForbidden, err.Error(), nil); writeErr != nil {
			l.Error("write list projects error response", "error", writeErr)
		}
		l.Error("get user ID for list projects", "error", err)
		return
	}

	projects, err := h.st.ListProjectsByUserID(r.Context(), userID)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, err.Error(), nil); writeErr != nil {
			l.Error("write list projects error response", "error", writeErr)
		}
		l.Error("list projects", "error", err)
		return
	}

	if err := HTTPResponse(w, http.StatusOK, "projects retrieved successfully", projects); err != nil {
		l.Error("write list projects response", "error", err)
	}
	l.Info("projects retrieved successfully", "user_id", userID.String())
}

func (h *ServerHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)

	projectIDStr := r.PathValue("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusBadRequest, "invalid project ID", nil); writeErr != nil {
			l.Error("write get project error response", "error", writeErr)
		}
		l.Error("parse project ID", "error", err)
		return
	}

	project, err := h.st.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, err.Error(), nil); writeErr != nil {
			l.Error("write get project error response", "error", writeErr)
		}
		l.Error("get project", "error", err)
		return
	}

	if err := HTTPResponse(w, http.StatusOK, "project retrieved successfully", project); err != nil {
		l.Error("write get project response", "error", err)
	}
	l.Info("project retrieved successfully", "project_id", projectID.String())
}

func (h *ServerHandler) UpdateProjectByID(w http.ResponseWriter, r *http.Request) {}

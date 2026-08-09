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
		HTTPError(w, http.StatusForbidden, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		HTTPError(w, http.StatusBadRequest, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	projectID, err := h.st.CreateProject(r.Context(), userID, req)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, "failed to create project", nil)
		l.Error(err.Error())
		return
	}

	HTTPResponse(w, http.StatusCreated, "project created successfully", map[string]string{"project_id": projectID.String()})
	l.Info("project created successfully", "user_id", userID.String(), "project_id", projectID.String())
}

func (h *ServerHandler) GetAllProjectsForUser(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	userID, err := m.GetUserIDFromContext(r.Context())
	if err != nil {
		HTTPError(w, http.StatusForbidden, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	projects, err := h.st.ListProjectsByUserID(r.Context(), userID)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	HTTPResponse(w, http.StatusOK, "projects retrieved successfully", projects)
	l.Info("projects retrieved successfully", "user_id", userID.String())
}

func (h *ServerHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)

	projectIDStr := r.PathValue("projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "invalid project ID", nil)
		l.Error(err.Error())
		return
	}

	project, err := h.st.GetProjectByID(r.Context(), projectID)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	HTTPResponse(w, http.StatusOK, "project retrieved successfully", project)
	l.Info("project retrieved successfully", "project_id", projectID.String())
}

func (h *ServerHandler) UpdateProjectByID(w http.ResponseWriter, r *http.Request) {}

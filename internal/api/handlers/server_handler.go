package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/Diaku49/AI-visibility-tracker/internal/logger"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
	"github.com/google/uuid"
)

type ServerHandler struct {
	store store.Store
	l     *slog.Logger
}

func NewServerHandler(store store.Store) *ServerHandler {
	return &ServerHandler{
		store: store,
		l:     logger.NewLogger(),
	}
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func HTTPResponse(w http.ResponseWriter, statusCode int, msg string, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{
		Message: msg,
		Data:    data,
	}); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func HTTPError(w http.ResponseWriter, statusCode int, errMsg string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(Response{
		Message: errMsg,
		Data:    data,
	}); err != nil {
		log.Printf("failed to write error JSON: %v", err)
	}
}

func (h *ServerHandler) GetLogger(r *http.Request) *slog.Logger {
	reqID := r.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = uuid.New().String()
	}

	log := h.l.With(
		"request_id", reqID,
		"path", r.URL.Path,
		"method", r.Method,
	)

	return log
}

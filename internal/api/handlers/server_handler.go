package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/Diaku49/AI-visibility-tracker/config"
	"github.com/Diaku49/AI-visibility-tracker/internal/logger"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ServerHandler struct {
	st        *store.Store
	l         *slog.Logger
	v         *validator.Validate
	cfg       *config.Config
	keyCipher *pkg.KeyCipher
}

func NewServerHandler(store *store.Store, cfg *config.Config, kc *pkg.KeyCipher) *ServerHandler {
	pkg.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &ServerHandler{
		st:        store,
		l:         logger.NewLogger(),
		v:         pkg.Validate,
		cfg:       cfg,
		keyCipher: kc,
	}
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func decodeAndValidateJSON(r *http.Request, v any, validate *validator.Validate) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("invalid json body: %w", err)
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain only one JSON object")
	}

	if err := validate.Struct(v); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func HTTPResponse(w http.ResponseWriter, statusCode int, msg string, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{
		Message: msg,
		Data:    data,
	}); err != nil {
		return fmt.Errorf("encode HTTP response: %w", err)
	}

	return nil
}

func HTTPError(w http.ResponseWriter, statusCode int, errMsg string, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(Response{
		Message: errMsg,
		Data:    data,
	}); err != nil {
		return fmt.Errorf("encode HTTP error response: %w", err)
	}

	return nil
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

package handlers

import (
	"net/http"

	"github.com/Diaku49/AI-visibility-tracker/internal/api/middleware"
	"github.com/Diaku49/AI-visibility-tracker/internal/dto"
)

func (h *ServerHandler) CreateProviderKey(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	var req dto.CreateProviderKeyRequest
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		if writeErr := HTTPError(w, http.StatusForbidden, err.Error(), nil); writeErr != nil {
			l.Error("write create provider key error response", "error", writeErr)
		}
		l.Error("get user ID for create provider key", "error", err)
		return
	}

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		if writeErr := HTTPError(w, http.StatusBadRequest, err.Error(), nil); writeErr != nil {
			l.Error("write create provider key error response", "error", writeErr)
		}
		l.Error("create provider key request validation failed", "error", err)
		return
	}

	encryptedKey, nounce, err := h.keyCipher.Encrypt([]byte(req.Key))
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, "failed to encrypt provider key", nil); writeErr != nil {
			l.Error("write create provider key error response", "error", writeErr)
		}
		l.Error("encrypt provider key", "error", err)
		return
	}

	ID, err := h.st.CreateProviderKey(
		r.Context(),
		userID,
		req.EngineID,
		req.Name,
		encryptedKey,
		nounce,
		true,
		nil,
	)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, "failed to create provider key", nil); writeErr != nil {
			l.Error("write create provider key error response", "error", writeErr)
		}
		l.Error("create provider key", "error", err)
		return
	}

	if err := HTTPResponse(w, http.StatusCreated, "key created successfully", nil); err != nil {
		l.Error("write create provider key response", "error", err)
	}
	l.Info("key created successfully", "keyID", ID)
}

func (h *ServerHandler) DeleteProviderKey(w http.ResponseWriter, r *http.Request) {

}

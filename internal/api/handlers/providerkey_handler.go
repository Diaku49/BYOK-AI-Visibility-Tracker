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
		HTTPError(w, http.StatusForbidden, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		HTTPError(w, http.StatusBadRequest, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	encryptedKey, nounce, err := h.keyCipher.Encrypt([]byte(req.Key))
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, "failed to encrypt provider key", nil)
		l.Error(err.Error())
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

	HTTPResponse(w, http.StatusCreated, "key created successfully", nil)
	l.Info("key created successfully", "keyID", ID)
}

func (h *ServerHandler) DeleteProviderKey(w http.ResponseWriter, r *http.Request) {

}

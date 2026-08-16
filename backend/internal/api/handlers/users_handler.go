package handlers

import (
	"net/http"

	"github.com/Diaku49/AI-visibility-tracker/internal/dto"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
)

func (h *ServerHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	var req dto.SignUpUser

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		if writeErr := HTTPError(w, http.StatusBadRequest, err.Error(), nil); writeErr != nil {
			l.Error("write sign-up error response", "error", writeErr)
		}
		l.Error("sign-up request validation failed", "error", err)
		return
	}

	userID, err := h.st.CreateUser(r.Context(), req.Email, req.Email, req.Name)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, err.Error(), nil); writeErr != nil {
			l.Error("write sign-up error response", "error", writeErr)
		}
		l.Error("create user", "error", err)
		return
	}

	if err := HTTPResponse(w, http.StatusCreated, "user signed up.", nil); err != nil {
		l.Error("write sign-up response", "error", err)
	}
	l.Info("signed up successfully", "userID", userID.String())
}

func (h *ServerHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	var req dto.LoginUser

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		if writeErr := HTTPError(w, http.StatusBadRequest, err.Error(), nil); writeErr != nil {
			l.Error("write login error response", "error", writeErr)
		}
		l.Error("login request validation failed", "error", err)
		return
	}

	userID, tier, pass, err := h.st.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, err.Error(), nil); writeErr != nil {
			l.Error("write login error response", "error", writeErr)
		}
		l.Error("get user for login", "error", err)
		return
	}

	if !pkg.CheckPasswordHash(req.Password, pass) {
		if writeErr := HTTPError(w, http.StatusForbidden, "wrong password", nil); writeErr != nil {
			l.Error("write login error response", "error", writeErr)
		}
		l.Error("login rejected: wrong password", "user_id", userID)
		return
	}

	jwt, err := pkg.GenerateJWT(userID, tier, []byte(h.cfg.JWTSecret))
	if err != nil {
		if writeErr := HTTPError(w, http.StatusInternalServerError, err.Error(), nil); writeErr != nil {
			l.Error("write login error response", "error", writeErr)
		}
		l.Error("generate login token", "error", err)
		return
	}

	if err := HTTPResponse(w, http.StatusOK, "logged in successfully", jwt); err != nil {
		l.Error("write login response", "error", err)
	}
	l.Info("loggedn in successfully", "userID", userID.String())
}

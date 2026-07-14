package handlers

import (
	"context"
	"net/http"

	"github.com/Diaku49/AI-visibility-tracker/internal/model"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
)

func (h *ServerHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	ctx := context.Background()
	var req model.SignUpUser

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		HTTPError(w, http.StatusBadRequest, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	userID, err := h.st.CreateUser(ctx, req.Email, req.Email, req.Name)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	HTTPResponse(w, http.StatusCreated, "user signed up.", nil)
	l.Info("signed up successfully", "userID", userID.String())
}

func (h *ServerHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	l := h.GetLogger(r)
	ctx := context.Background()
	var req model.LoginUser

	if err := decodeAndValidateJSON(r, &req, h.v); err != nil {
		HTTPError(w, http.StatusBadRequest, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	userID, tier, pass, err := h.st.GetUserByEmail(ctx, req.Email)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	if !pkg.CheckPasswordHash(req.Password, pass) {
		HTTPError(w, http.StatusForbidden, "wrong password", nil)
		l.Error("wrong password")
		return
	}

	jwt, err := pkg.GenerateJWT(userID, tier, []byte(h.cfg.JWTSecret))
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err.Error(), nil)
		l.Error(err.Error())
		return
	}

	HTTPResponse(w, http.StatusOK, "logged in successfully", jwt)
	l.Info("loggedn in successfully", "userID", userID.String())
}

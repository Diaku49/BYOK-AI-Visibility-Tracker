package api

import (
	"net/http"

	"github.com/Diaku49/AI-visibility-tracker/config"
	"github.com/Diaku49/AI-visibility-tracker/internal/api/handlers"
	"github.com/Diaku49/AI-visibility-tracker/internal/api/middleware"
	"github.com/Diaku49/AI-visibility-tracker/internal/pkg"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
)

type Server struct {
	store store.Store
	cfg   *config.Config
	kc    *pkg.KeyCipher
}

func NewServer(store store.Store, config *config.Config) *Server {
	return &Server{
		store: store,
		cfg:   config,
	}
}

func (s *Server) Route(h *handlers.ServerHandler) http.Handler {
	mux := http.NewServeMux()
	authMiddleware := middleware.Authenticate(s.cfg.JWTSecret)

	// users
	mux.HandleFunc("POST /user", h.SignUpUser)
	mux.HandleFunc("POST /user/login", h.LoginUser)

	// keys
	mux.Handle("POST /key", authMiddleware(http.HandlerFunc(h.CreateProviderKey)))

	// projects

	// subscriptions

	// runs

	return mux
}

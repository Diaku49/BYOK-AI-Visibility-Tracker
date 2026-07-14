package api

import (
	"net/http"

	"github.com/Diaku49/AI-visibility-tracker/config"
	"github.com/Diaku49/AI-visibility-tracker/internal/api/handlers"
	"github.com/Diaku49/AI-visibility-tracker/internal/store"
)

type Server struct {
	store store.Store
	cfg   *config.Config
}

func NewServer(store store.Store, config *config.Config) *Server {
	return &Server{
		store: store,
		cfg:   config,
	}
}

func (s *Server) Route() http.Handler {
	mux := http.NewServeMux()
	h := handlers.NewServerHandler(s.store)
	authMiddleware := Authenticate(s.cfg.JWTSecret)

	// users
	mux.HandleFunc("POST /user", h.SignUpUser)

	// keys
	mux.Handle("POST /key", authMiddleware(http.HandlerFunc(h.CreateKey)))

	// projects

	// subscriptions

	// runs

	return mux
}

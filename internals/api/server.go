package api

import (
	"fmt"
	"net/http"

	"github.com/Diaku49/AI-visibility-tracker/internals/api/handlers"
	"github.com/Diaku49/AI-visibility-tracker/internals/store"
)

type Server struct {
	store store.Store
}

func NewServer(store store.Store) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) Route() http.Handler {
	mux := http.NewServeMux()
	h := handlers.NewServerHandler(s.store)
	fmt.Printf("%v", h)

	// users

	// projects

	// subscriptions

	// runs

	return mux
}

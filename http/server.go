package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// Server is a wrapper around the http server and router and and dependencies the handlers may need
type Server struct {
	router *chi.Mux

	Addr string
	Port string
}

// NewServer creates a new server instance
func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.setupRoutes()

	return s
}

// Open starts listening for requests
func (s *Server) Open() error {
	addr := fmt.Sprintf("%s:%s", s.Addr, s.Port)

	zap.S().Infof("Starting API http://%s", addr)

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

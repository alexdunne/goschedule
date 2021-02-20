package http

import "net/http"

func (s *Server) setupRoutes() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	s.router.Post("/auth/login", s.handleAuthLogin)
}

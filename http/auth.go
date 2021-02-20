package http

import (
	"goschedule/internal/render"
	"net/http"
)

func (s *Server) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.H{
		"user": render.H{
			"id": "abc",
		},
	})
}

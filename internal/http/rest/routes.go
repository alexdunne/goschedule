package rest

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"goschedule/internal/accounts"
	"goschedule/internal/render"
	"io"
	"net/http"
	"strconv"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func (s *Server) setupRoutes() {
	s.router.Get("/oauth/github", s.handleOAuthGithub)
	s.router.Get("/oauth/github/callback", s.handleOAuthGitHubCallback)
}

func (s *Server) handleOAuthGithub(w http.ResponseWriter, r *http.Request) {
	session, err := s.session(r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	state := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, state); err != nil {
		render.Error(w, r, err)
		return
	}
	session.State = hex.EncodeToString(state)

	if err := s.setSession(w, r, session); err != nil {
		render.Error(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.H{
		"data": render.H{
			"url": s.GithubOAuth2Config().AuthCodeURL(session.State),
		},
	})
}

func (s *Server) handleOAuthGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state, code := r.FormValue("state"), r.FormValue("code")

	session, err := s.session(r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if state != session.State {
		render.Error(w, r, fmt.Errorf("oauth state mismatch"))
		return
	}

	tok, err := s.GithubOAuth2Config().Exchange(r.Context(), code)
	if err != nil {
		render.Error(w, r, fmt.Errorf("oauth exchange error: %s", err))
		return
	}

	client := github.NewClient(oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tok.AccessToken},
	)))

	u, _, err := client.Users.Get(r.Context(), "")
	if err != nil {
		render.Error(w, r, fmt.Errorf("cannot fetch github user: %s", err))
		return
	} else if u.ID == nil {
		render.Error(w, r, fmt.Errorf("user ID not returned by GitHub, cannot authenticate user"))
		return
	}

	var name string
	if u.Name != nil {
		name = *u.Name
	} else if u.Login != nil {
		name = *u.Login
	}
	var email string
	if u.Email != nil {
		email = *u.Email
	}

	account := &accounts.NewAccount{
		Name:     name,
		Email:    email,
		Source:   accounts.AuthSourceGitHub,
		SourceID: strconv.FormatInt(*u.ID, 10),
	}

	if err := s.AccountsService.CreateAccount(r.Context(), account); err != nil {
		render.Error(w, r, err)
		return
	}

	session.UserID = account.UserID
	session.State = ""
	if err := s.setSession(w, r, session); err != nil {
		render.Error(w, r, fmt.Errorf("cannot set session cookie: %s", err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.H{
		"data": render.H{
			"user": render.H{
				"id": account.UserID,
			},
		},
	})
}

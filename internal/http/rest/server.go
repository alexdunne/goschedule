package rest

import (
	"fmt"
	"goschedule/internal/accounts"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Server is a wrapper around the http server and router and and dependencies the handlers may need
type Server struct {
	router       *chi.Mux
	sessionStore *sessions.CookieStore

	Addr string
	Port string

	// Keys used for secure cookie encryption
	HashKey  string
	BlockKey string

	// Credentials for Github OAuth
	GitHubClientID     string
	GitHubClientSecret string

	// Services
	AccountsService accounts.Service
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
	if s.Addr == "" {
		return fmt.Errorf("addr is required")
	} else if s.Port == "" {
		return fmt.Errorf("port is required")
	} else if s.GitHubClientID == "" {
		return fmt.Errorf("github client id required")
	} else if s.GitHubClientSecret == "" {
		return fmt.Errorf("github client secret required")
	}

	if err := s.enableSessionStore(); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%s", s.Addr, s.Port)

	zap.S().Infof("Starting API http://%s", addr)

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) enableSessionStore() error {
	zap.S().Info("Setting up session store")

	if s.HashKey == "" {
		return fmt.Errorf("hash key required")
	} else if s.BlockKey == "" {
		return fmt.Errorf("block key required")
	}

	s.sessionStore = sessions.NewCookieStore([]byte(s.HashKey), []byte(s.BlockKey))
	s.sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	return nil
}

// SessionCookieName is the name of the cookie used to store the session
const SessionCookieName = "__goschedule_session__"

// Session represents session data stored in a secure cookie
type Session struct {
	UserID string `json:"userID"`
	State  string `json:"state"`
}

func (s *Server) session(r *http.Request) (Session, error) {
	session, _ := s.sessionStore.Get(r, SessionCookieName)

	if session.IsNew {
		return Session{}, nil
	}

	return Session{
		UserID: session.Values["UserID"].(string),
		State:  session.Values["State"].(string),
	}, nil
}

func (s *Server) setSession(w http.ResponseWriter, r *http.Request, session Session) error {
	sess, _ := s.sessionStore.Get(r, SessionCookieName)

	sess.Values["UserID"] = session.UserID
	sess.Values["State"] = session.State

	if err := sess.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *Server) clearSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := s.sessionStore.Get(r, SessionCookieName)

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}

func (s *Server) GithubOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.GitHubClientID,
		ClientSecret: s.GitHubClientSecret,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}

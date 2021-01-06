package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/zerok/retoots/pkg/mastodon"
)

type Server struct {
	r                   chi.Router
	mc                  *mastodon.Client
	allowedRootAccounts []string
}

type Configurator func(*Config)

type Config struct {
	MastodonClient      *mastodon.Client
	AllowedOrigins      []string
	AllowedRootAccounts []string
}

func WithMastodonClient(mc *mastodon.Client) Configurator {
	return func(cfg *Config) {
		cfg.MastodonClient = mc
	}
}

func WithAllowedOrigins(origins []string) Configurator {
	return func(cfg *Config) {
		cfg.AllowedOrigins = origins
	}
}

func WithAllowedRootAccounts(accounts []string) Configurator {
	return func(cfg *Config) {
		cfg.AllowedRootAccounts = accounts
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *Server) isFromAllowedRootAccount(ctx context.Context, u string) bool {
	fmt.Sprintf("[%s]\n", s.allowedRootAccounts)
	if s.allowedRootAccounts == nil || len(s.allowedRootAccounts) == 0 {
		return true
	}
	status, err := s.mc.GetStatus(ctx, u)
	if err != nil {
		return false
	}
	for _, a := range s.allowedRootAccounts {
		if a == status.Account.Acct {
			return true
		}
	}
	return false
}

func (s *Server) handleGetDescendants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	statusURL := r.URL.Query().Get("status")
	if statusURL == "" {
		http.Error(w, "no status URL", http.StatusBadRequest)
		return
	}
	if !s.isFromAllowedRootAccount(ctx, statusURL) {
		http.Error(w, "not allowed URL", http.StatusBadRequest)
		return
	}
	descendants, err := s.mc.GetDescendants(ctx, statusURL)
	if err != nil {
		http.Error(w, "failed", http.StatusBadRequest)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(descendants)
}

func New(ctx context.Context, configurators ...Configurator) *Server {
	cfg := Config{}
	for _, configurator := range configurators {
		configurator(&cfg)
	}
	cors := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowCredentials: true,
	})
	router := chi.NewRouter()
	router.Use(cors.Handler)
	srv := &Server{
		r:                   router,
		mc:                  cfg.MastodonClient,
		allowedRootAccounts: cfg.AllowedRootAccounts,
	}
	router.Get("/api/v1/descendants", srv.handleGetDescendants)
	return srv
}

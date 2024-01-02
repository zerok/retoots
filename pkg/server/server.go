package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/zerok/retoots/pkg/mastodon"
)

type Server struct {
	r                   chi.Router
	mc                  *mastodon.Client
	allowedRootAccounts []string
	rootCheckCache      *lru.Cache[string, bool]
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
	if s.allowedRootAccounts == nil || len(s.allowedRootAccounts) == 0 {
		return true
	}
	if result, fromCache := s.rootCheckCache.Get(u); fromCache {
		return result
	}
	status, err := s.mc.GetStatus(ctx, u)
	if err != nil {
		s.rootCheckCache.Add(u, false)
		return false
	}
	for _, a := range s.allowedRootAccounts {
		if a == status.Account.Acct {
			s.rootCheckCache.Add(u, true)
			return true
		}
	}
	s.rootCheckCache.Add(u, false)
	return false
}

func (s *Server) handleGetDescendants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	statusURL := r.URL.Query().Get("status")
	if statusURL == "" {
		http.Error(w, "no status URL", http.StatusBadRequest)
		return
	}
	statusURL, err := normalizeStatusURL(statusURL)
	if err != nil {
		http.Error(w, "not allowed URL", http.StatusBadRequest)
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

func (s *Server) handleGetFavoritedBy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	statusURL := r.URL.Query().Get("status")
	if statusURL == "" {
		http.Error(w, "no status URL", http.StatusBadRequest)
		return
	}
	statusURL, err := normalizeStatusURL(statusURL)
	if err != nil {
		http.Error(w, "not allowed URL", http.StatusBadRequest)
		return
	}
	if !s.isFromAllowedRootAccount(ctx, statusURL) {
		http.Error(w, "not allowed URL", http.StatusBadRequest)
		return
	}
	favorites, err := s.mc.GetFavoritedBy(ctx, statusURL)
	if err != nil {
		http.Error(w, "failed", http.StatusBadRequest)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

func (s *Server) handleInteractions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	statusURL := r.URL.Query().Get("status")
	if statusURL == "" {
		http.Error(w, "no status URL", http.StatusBadRequest)
		return
	}
	statusURL, err := normalizeStatusURL(statusURL)
	if err != nil {
		http.Error(w, "not allowed URL", http.StatusBadRequest)
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
	favorites, err := s.mc.GetFavoritedBy(ctx, statusURL)
	if err != nil {
		http.Error(w, "failed", http.StatusBadRequest)
	}
	boosts, err := s.mc.GetBoostedBy(ctx, statusURL)
	if err != nil {
		http.Error(w, "failed", http.StatusBadRequest)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interactions{
		Descendants: descendants,
		FavoritedBy: favorites,
		BoostedBy:   boosts,
	})
}

type interactions struct {
	Descendants []mastodon.Status  `json:"descendants"`
	FavoritedBy []mastodon.Account `json:"favorited_by"`
	BoostedBy   []mastodon.Account `json:"boosted_by"`
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

	// lru.New can only fail if a negative size is provided.
	cache, _ := lru.New[string, bool](1024)
	srv := &Server{
		r:                   router,
		mc:                  cfg.MastodonClient,
		allowedRootAccounts: cfg.AllowedRootAccounts,
		rootCheckCache:      cache,
	}
	router.Get("/api/v1/descendants", srv.handleGetDescendants)
	router.Get("/api/v1/favorited_by", srv.handleGetFavoritedBy)
	router.Get("/api/v1/interactions", srv.handleInteractions)
	return srv
}

func normalizeStatusURL(statusURL string) (string, error) {
	u, err := url.Parse(statusURL)
	if err != nil {
		return "", err
	}
	u.RawQuery = ""
	pathElements := strings.Split(u.Path, "/")
	if len(pathElements) > 3 {
		pathElements = pathElements[:3]
	}
	u.Path = strings.Join(pathElements, "/") + "/"
	return u.String(), nil
}

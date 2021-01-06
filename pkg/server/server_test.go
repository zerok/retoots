package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"github.com/zerok/retoots/pkg/mastodon"
)

func generateStatus(t *testing.T, id string, content string, authorID string, authorName string) string {
	return fmt.Sprintf(`{"id": "%s", "content": "%s", "account":{"id": "%s", "acct": "%s"}}`, id, content, authorID, authorName)
}

func TestGetDescendants(t *testing.T) {
	ctx := context.Background()
	router := chi.NewRouter()
	status123 := generateStatus(t, "123", "content123", "123", "username123")
	status234 := generateStatus(t, "234", "content234", "234", "username234")
	router.Get("/api/v1/statuses/123", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, status123)
	})
	router.Get("/api/v1/statuses/123/context", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{}`)
	})
	router.Get("/api/v1/statuses/234", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, status234)
	})
	router.Get("/api/v1/statuses/234/context", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{}`)
	})
	msrv := httptest.NewServer(router)

	t.Run("not-exists", func(t *testing.T) {
		surl := fmt.Sprintf("%s/@username123/000", msrv.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/descendants?status="+surl, nil)
		srv := New(ctx)
		srv.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("no-allowed-accounts", func(t *testing.T) {
		surl := fmt.Sprintf("%s/@username123/123", msrv.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/descendants?status="+surl, nil)
		srv := New(ctx)
		srv.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not-allowed-accounts", func(t *testing.T) {
		surl := fmt.Sprintf("%s/@username123/123", msrv.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/descendants?status="+surl, nil)
		srv := New(ctx, WithAllowedRootAccounts([]string{"username234"}))
		srv.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("allowed-accounts", func(t *testing.T) {
		surl := fmt.Sprintf("%s/@username123/123", msrv.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/v1/descendants?status="+surl, nil)
		srv := New(ctx, WithAllowedRootAccounts([]string{mastodon.NormalizeAcct(msrv.URL, "username123")}))
		srv.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})
}

package mastodon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHostIDSplitting(t *testing.T) {
	u, err := ParseStatusURL("https://chaos.social/@zerok/105475673515689340")
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.Server, "https://chaos.social")
	require.Equal(t, u.ID, "105475673515689340")
}

func TestGetDescendants(t *testing.T) {
	c := New()
	t.Run("empty-response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		requestURL := fmt.Sprintf("%s/api/v1/statuses/123/context", srv.URL)
		ctx := context.Background()
		_, err := c.GetDescendants(ctx, requestURL)
		require.Error(t, err)
	})
	t.Run("not-found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		}))
		requestURL := fmt.Sprintf("%s/api/v1/statuses/123/context", srv.URL)
		ctx := context.Background()
		_, err := c.GetDescendants(ctx, requestURL)
		require.Error(t, err)
	})
	t.Run("not-descendants", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
"descendants": []
}`)
		}))
		requestURL := fmt.Sprintf("%s/api/v1/statuses/123/context", srv.URL)
		ctx := context.Background()
		s, err := c.GetDescendants(ctx, requestURL)
		require.NoError(t, err)
		require.Len(t, s, 0)
	})
	t.Run("non-empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
"descendants": [{
  "id": "123",
  "content": "<p>hello</p>"
}]
}`)
		}))
		requestURL := fmt.Sprintf("%s/api/v1/statuses/123/context", srv.URL)
		ctx := context.Background()
		s, err := c.GetDescendants(ctx, requestURL)
		require.NoError(t, err)
		require.Len(t, s, 1)
		require.Equal(t, "123", s[0].ID)
		require.Equal(t, "<p>hello</p>", s[0].Content)
	})
}

func TestGetFavoritedBy(t *testing.T) {
	c := New()
	t.Run("non-empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `[{
  "id": "123",
  "username": "username",
  "acct": "username"
}]`)
		}))
		requestURL := fmt.Sprintf("%s/api/v1/statuses/123/favorited_by", srv.URL)
		u, _ := url.Parse(srv.URL)
		ctx := context.Background()
		accounts, err := c.GetFavoritedBy(ctx, requestURL)
		require.NoError(t, err)
		require.Len(t, accounts, 1)
		require.Equal(t, fmt.Sprintf("username@%s", u.Host), accounts[0].Acct)
	})
}

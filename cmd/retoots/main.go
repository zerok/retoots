package main

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/zerok/retoots/pkg/mastodon"
	"github.com/zerok/retoots/pkg/server"
)

func main() {
	var addr string
	var allowedOrigins []string
	var allowedRootAccounts []string
	pflag.StringVar(&addr, "addr", "localhost:8000", "Address to listen on for requests")
	pflag.StringSliceVar(&allowedOrigins, "allowed-origins", []string{"localhost:8000"}, "CORS allowed origins")
	pflag.StringSliceVar(&allowedRootAccounts, "allowed-root-accounts", []string{}, "Accounts that are allowed as root of discussions")
	pflag.Parse()
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	ctx := logger.WithContext(context.Background())
	srv := http.Server{}
	srv.Addr = addr
	srv.Handler = server.New(ctx,
		server.WithMastodonClient(mastodon.New()),
		server.WithAllowedRootAccounts(allowedRootAccounts),
		server.WithAllowedOrigins(allowedOrigins))
	logger.Info().Msgf("Starting server on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("Server exited")
	}
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/zerok/retoots/pkg/mastodon"
	"github.com/zerok/retoots/pkg/server"
)

var version, commit, date string

func main() {
	var addr string
	var allowedOrigins []string
	var allowedRootAccounts []string
	var showVersion bool
	pflag.StringVar(&addr, "addr", "localhost:8000", "Address to listen on for requests")
	pflag.StringSliceVar(&allowedOrigins, "allowed-origins", []string{"localhost:8000"}, "CORS allowed origins")
	pflag.StringSliceVar(&allowedRootAccounts, "allowed-root-accounts", []string{}, "Accounts that are allowed as root of discussions")
	pflag.BoolVar(&showVersion, "version", false, "Show version information")
	pflag.Parse()
	if showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nBuild date: %s\n", version, commit, date)
		os.Exit(0)
	}
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

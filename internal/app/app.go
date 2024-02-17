package app

import (
	"log/slog"

	"github.com/ysomad/golangconf-bot/internal/config"
	"github.com/ysomad/golangconf-bot/internal/postgres"
	"github.com/ysomad/golangconf-bot/internal/slogx"
)

func Run(conf *config.Config) {
	slog.Debug("starting app", "config", conf)

	_, err := postgres.NewClient(conf.PG.URL, postgres.WithMaxConns(conf.PG.MaxConns))
	if err != nil {
		slogx.Fatal("postgres client not created", err)
	}
}

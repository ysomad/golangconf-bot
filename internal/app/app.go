package app

import (
	"log/slog"

	"github.com/ysomad/golangconf-bot/internal/config"
)

func Run(conf *config.Config) {
	slog.Debug("starting app", "config", conf)
}

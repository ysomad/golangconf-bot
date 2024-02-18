package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jackc/pgx/v5/pgxpool"
	tele "gopkg.in/telebot.v3"

	"github.com/ysomad/golangconf-bot/internal/bot"
	"github.com/ysomad/golangconf-bot/internal/config"
	"github.com/ysomad/golangconf-bot/internal/slogx"
	"github.com/ysomad/golangconf-bot/internal/state"
)

func Run(conf *config.Config) {
	slog.Debug("starting app", "config", conf)

	// postgres
	pgPoolConf, err := pgxpool.ParseConfig(conf.PG.URL)
	if err != nil {
		slogx.Fatal("postgres connection url not parsed", err)
	}

	pgPoolConf.MaxConns = conf.PG.MaxConns

	pgPool, err := pgxpool.NewWithConfig(context.Background(), pgPoolConf)
	if err != nil {
		slogx.Fatal("postgres pool not created", err)
	}
	defer pgPool.Close()

	connAttempts := 10

	for connAttempts > 0 {
		if err = pgPool.Ping(context.Background()); err == nil {
			break
		}

		slog.Info("trying connecting to postgres", "attempts_left", connAttempts)
		time.Sleep(time.Second)
		connAttempts--
	}

	if err != nil {
		slogx.Fatal("postgres not available", err)
	}

	// sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// lru
	stateLRU := expirable.NewLRU[string, state.State](conf.StateLRU.Size, func(k string, v state.State) {
		slog.Warn("key evicted from state lru, consider increasing cache size", "key", k, "value", v)
	}, conf.StateLRU.TTL)

	stateStorage := state.NewStorage(stateLRU)

	// telegram
	telebot, err := tele.NewBot(tele.Settings{
		Token:   conf.Telegram.Token,
		Verbose: conf.Telegram.Verbose,
		OnError: bot.HandleError,
		Client:  &http.Client{Timeout: conf.Telegram.HTTPTimeout},

		// TODO: CHANGE TO WEBHOOK
		Poller: &tele.LongPoller{Timeout: time.Second},
	})
	if err != nil {
		slogx.Fatal("telebot not created", err)
	}

	appbot, err := bot.New(telebot, conf.Telegram.Admins, stateStorage)
	if err != nil {
		slogx.Fatal("app bot not created", err)
	}

	appbot.Start()
}

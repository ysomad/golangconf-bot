package bot

import (
	"log/slog"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v3"
	telemiddleware "gopkg.in/telebot.v3/middleware"

	"github.com/ysomad/golangconf-bot/internal/config"
)

type Bot struct {
	*tele.Bot
}

func New(conf config.Telegram) (Bot, error) {
	telebot, err := tele.NewBot(tele.Settings{
		Token:   conf.Token,
		Verbose: conf.Verbose,
		OnError: handleError,
		Client:  &http.Client{Timeout: conf.HTTPTimeout},

		// TODO: CHANGE TO WEBHOOK
		Poller: &tele.LongPoller{Timeout: time.Second},
	})
	if err != nil {
		return Bot{}, err
	}

	b := Bot{telebot}

	b.Use(telemiddleware.Recover())

	b.Handle(tele.OnText, b.HandleText)
	b.Handle(tele.OnCallback, b.HandleCallback)

	adminOnly := b.Group()
	adminOnly.Use(middlewareAdminOnly(conf.Admins))

	return b, nil
}

func handleError(err error, c tele.Context) {
	if err == nil || c == nil {
		return
	}

	slog.Error("shit happened", "error", err.Error(), "tg_id", c.Chat().ID)
}

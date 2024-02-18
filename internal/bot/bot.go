package bot

import (
	"log/slog"

	"github.com/ysomad/golangconf-bot/internal/state"
	tele "gopkg.in/telebot.v3"
	telemiddleware "gopkg.in/telebot.v3/middleware"
)

type Bot struct {
	*tele.Bot
	state stateStorage
}

type stateStorage interface {
	Save(int64, state.State)
	Get(int64) state.State
	Del(int64)
}

func New(telebot *tele.Bot, admins []int64, state stateStorage) (Bot, error) {
	b := Bot{
		Bot:   telebot,
		state: state,
	}

	b.Use(telemiddleware.Recover())

	b.Handle(tele.OnText, b.handleText)
	b.Handle(tele.OnCallback, b.handleCallback)

	adminOnly := b.Group()
	adminOnly.Use(middlewareAdminOnly(admins))

	adminOnly.Handle("/upload_schedule", b.uploadSchedule)
	adminOnly.Handle(tele.OnDocument, b.handleDocumentUpload)

	return b, nil
}

func HandleError(err error, c tele.Context) {
	if err == nil || c == nil {
		return
	}

	slog.Error("shit happened", "error", err.Error(), "tg_id", c.Chat().ID)
}

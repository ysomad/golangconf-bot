package bot

import (
	"github.com/ysomad/golangconf-bot/internal/state"
	"gopkg.in/telebot.v3"
)

func (b Bot) uploadSchedule(c telebot.Context) error {
	b.state.Save(c.Chat().ID, state.State{Step: state.StepUploadingSchedule})
	return c.Send("Отправь мне файл с расписанием докладов")
}

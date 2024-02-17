package bot

import "gopkg.in/telebot.v3"

func (b Bot) HandleCallback(c telebot.Context) error {
	return c.Send("CALLBACK NOT IMPLEMENTED")
}

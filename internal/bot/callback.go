package bot

import "gopkg.in/telebot.v3"

func (b Bot) handleCallback(c telebot.Context) error {
	return c.Send("CALLBACK NOT IMPLEMENTED")
}

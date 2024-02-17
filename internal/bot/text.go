package bot

import tele "gopkg.in/telebot.v3"

func (b Bot) HandleText(c tele.Context) error {
	return c.Send("HANDLE TEXT NOT IMPLEMENTED")
}

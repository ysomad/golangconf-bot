package bot

import (
	"errors"
	"slices"

	tele "gopkg.in/telebot.v3"
)

func middlewareAdminOnly(admins []int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !slices.Contains(admins, c.Chat().ID) {
				return errors.New("access denied")
			}

			return next(c)
		}
	}
}

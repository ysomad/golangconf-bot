package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ysomad/golangconf-bot/internal/schedule"
	"github.com/ysomad/golangconf-bot/internal/state"
	"gopkg.in/telebot.v3"
)

func (b Bot) handleDocumentUpload(c telebot.Context) error {
	usr := c.Sender()
	if usr == nil {
		return errors.New("user not present in telebot context")
	}

	msg := c.Message()
	if msg == nil {
		return errors.New("message not present in telebot context")
	}

	st := b.state.Get(usr.ID)

	switch st.Step {
	case state.StepUploadingSchedule:
		if msg.Document == nil {
			return c.Send("Файл не найден в сообщении")
		}

		if !strings.HasSuffix(msg.Document.FileName, "csv") || msg.Document.MIME != "text/csv" {
			return c.Send("Файл расписания должен быть в CSV формате")
		}

		// file obj from telegram api
		tgfile := msg.Document.MediaFile()

		slog.Debug("got schedule file", "file", tgfile)

		if !tgfile.InCloud() {
			return c.Send("Файл расписания не найден на серверах Telegram, загрузи еще раз")
		}

		// downloaded file
		file, err := b.File(tgfile)
		if err != nil {
			return c.Send(fmt.Sprintf("Не удалось скачать файл с серверов Telegram, ошибка: %s", err.Error()))
		}
		defer file.Close()

		schedule, errs := schedule.NewFromCSV(file)
		if len(errs) != 0 {
			for _, err := range errs {
				c.Send(err.Error())
			}

			return nil
		}

		bb, err := json.Marshal(schedule)
		if err != nil {
			c.Send(err.Error())
		}

		return c.Send(string(bb))

	default:
		return fmt.Errorf("unsupported document upload state step: %d", st.Step)
	}
}

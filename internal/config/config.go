package config

import (
	"errors"
	"time"
)

type Config struct {
	Log        Log
	Telegram   Telegram
	PG         PG
	Conference Conference
}

type Log struct {
	Level     string `env:"LOG_LEVEL" env-default:"info"`
	AddSource bool   `env:"LOG_ADD_SOURCE"`
}

type Telegram struct {
	Token       string        `env:"TG_TOKEN" env-required:"true"`
	Admins      []int64       `env:"TG_ADMINS" env-required:"true"`
	HTTPTimeout time.Duration `env:"TG_HTTP_TIMEOUT" env-required:"true"`
	Verbose     bool          `env:"TG_DEBUG"`
}

type PG struct {
	URL      string `env:"PG_URL" env-required:"true"`
	MaxConns int32  `env:"PG_MAX_CONNS" env-required:"true"`
}

// Conference represents conferences that is taking place at the current moment in time.
// Move to MODELS ???
type Conference struct {
	Name               string    `env:"CONF_NAME" env-required:"true"`
	URL                string    `env:"CONF_URL" env-required:"true"`
	DateFrom           time.Time `env:"DATE_FROM" env-required:"true" env-layout:"02.01.2006"`
	DateUntil          time.Time `env:"DATE_UNTIL" env-required:"true" env-layout:"02.01.2006"`
	EvaluationDeadline time.Time `env:"EVALUATION_DEADLINE" env-required:"true" env-layout:"02.01.2006"`
}

func (c Conference) Validate() error {
	if c.DateUntil.Compare(c.DateFrom) != 1 {
		return errors.New("DATE_UNTIL must be greater than DATE_FROM")
	}

	if c.EvaluationDeadline.Compare(c.DateUntil) != 1 {
		return errors.New("EVALUATION_DEADLINE must be greater than DATE_UNTIL")
	}

	return nil
}

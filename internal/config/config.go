package config

import (
	"errors"
	"time"
)

type Config struct {
	PGURL      string  `env:"PG_URL" env-required:"true"`
	TGToken    string  `env:"TG_TOKEN" env-required:"true"`
	TGAdmins   []int64 `env:"TG_ADMINS" env-required:"true"`
	Conference Conference
	Log        Log
}

type Log struct {
	Level     string `env:"LOG_LEVEL" env-default:"info"`
	AddSource bool   `env:"LOG_ADD_SOURCE"`
}

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

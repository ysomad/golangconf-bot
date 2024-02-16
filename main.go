package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/ysomad/golangconf-bot/internal/app"
	"github.com/ysomad/golangconf-bot/internal/config"
	"github.com/ysomad/golangconf-bot/internal/slogx"
)

func main() {
	var (
		conf   config.Config
		logFmt string
	)

	flag.StringVar(&logFmt, "log-format", "text", "logging format, options: text, json")
	flag.Parse()

	if err := cleanenv.ReadEnv(&conf); err != nil {
		slogx.Fatal("config not parsed", err)
	}

	slogger := slog.New(slogx.NewHandler(os.Stdout, logFmt, conf.Log.Level, conf.Log.AddSource))
	slog.SetDefault(slogger)

	if err := conf.Conference.Validate(); err != nil {
		slogx.Fatal("config not valid", err)
	}

	app.Run(&conf)
}

package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v11"
)

// Config
// Основная конфигурационная структура, поле Config.Address отвечают за адрес запуска приложения
// Config.DatabaseUri за строку с данными для подключения к базе
// Config.AccrualSystemAddress адрес системы расчёта начислений
// /**
type Config struct {
	Address              string `env:"RUN_ADDRESS"`
	DatabaseUri          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCURAL_SYSTEM_ADDRESS"`
}

// ParseAppOptions парсит входящие параметры энва, и входящие аргументы запустившие процесс
func ParseAppOptions(environment []string, arguments []string) (*Config, error) {
	cfg := &Config{
		Address:              ":8080",
		DatabaseUri:          "user=user password=password dbname=pg sslmode=disable host=localhost port=5433",
		AccrualSystemAddress: "localhost",
	}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&cfg.Address, "a", cfg.Address, "run server address")
	fs.StringVar(&cfg.DatabaseUri, "d", cfg.DatabaseUri, "")
	fs.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "")

	err := fs.Parse(arguments)
	if err != nil {
		return nil, err
	}

	err = env.ParseWithOptions(cfg, env.Options{
		Environment: toMap(environment),
	})
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func toMap(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		if len(p) == 2 {
			r[p[0]] = p[1]
		}
	}
	return r
}

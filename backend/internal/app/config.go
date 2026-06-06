package app

import "log/slog"

type Config struct {
	LogLevel                 slog.Level `envconfig:"LOG_LEVEL"`
	Port                     int        `envconfig:"PORT" required:"true"`
	ExternalWebServerAddress string     `envconfig:"EXTERNAL_WEB_SERVER_ADDRESS"`
}

package config

import (
	"time"

	"github.com/rs/zerolog"
)

type OpenWeather struct {
	APIID   string        `required:"true"`
	BaseURL string        `required:"true"`
	Timeout time.Duration `default:"5s"`
}

type AuthService struct {
	URL string `default:"http://some.auth.com"`
}

type Config struct {
	Address          string        `default:"0.0.0.0"`
	Port             int           `default:"80"`
	LogLevel         zerolog.Level `default:"info"`
	ReadWriteTimeout time.Duration `default:"20s"`
	IdleTimeout      time.Duration `default:"75s"`
	ShutdownTime     time.Duration `default:"20s"`
	OpenWeather      OpenWeather
	AuthService      AuthService
}

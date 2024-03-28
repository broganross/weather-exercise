package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/broganross/weather-exercise/internal/config"
	"github.com/broganross/weather-exercise/internal/domain"
	"github.com/broganross/weather-exercise/internal/repo"
	"github.com/broganross/weather-exercise/internal/server"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "development"

func main() {
	// Load config from env vars
	conf := config.Config{}
	if err := envconfig.Process("weather", &conf); err != nil {
		log.Error().Err(err).Msg("parsing environment variables")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(conf.LogLevel)

	// construct services
	openWeather := &repo.OpenWeather{
		BaseURL: conf.OpenWeather.BaseURL,
		Client:  &http.Client{},
		APIid:   conf.OpenWeather.APIID,
		Timeout: conf.OpenWeather.Timeout,
	}
	domainService := &domain.WeatherService{
		Source: openWeather,
	}
	handlers := server.Handlers{
		Domain: domainService,
	}
	router := mux.NewRouter()
	server.SetupRoutes(&handlers, router, &conf)

	// create the server, and start it up
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		WriteTimeout: conf.ReadWriteTimeout,
		IdleTimeout:  conf.IdleTimeout,
		Handler:      router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("serving")
		}
	}()

	// Graceful shutdown
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTime)
	defer cancel()
	srv.Shutdown(ctx)
	log.Info().Msg("shutting down")
	os.Exit(1)

}

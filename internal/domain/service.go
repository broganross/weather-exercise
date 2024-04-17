package domain

import (
	"context"
	"fmt"

	"github.com/broganross/weather-exercise/internal/types"
)

type Temperature string

const (
	TempUnknown Temperature = "unknown"
	TempHot     Temperature = "hot"
	TempCold    Temperature = "cold"
	TempMod     Temperature = "moderate"
)

type Weather struct {
	Coords      types.Coords
	States      []string
	Temperature Temperature
}

// RepoWeather purely existing so that WeatherService.CurrentIn actually does something.
// In a normal case we would convert the repo data into domain data.  AKA join states, and convert the temperature.
type RepoWeather struct {
	Coords      types.Coords
	States      []string
	Temperature float32
}

// Exported Business logic interface
type Service interface {
	CurrentIn(ctx context.Context, lat float32, lon float32) (*Weather, error)
}

// Interface for where we're getting actual weather data from
type Repo interface {
	GetByCoords(ctx context.Context, latitude float32, longitude float32) (*RepoWeather, error)
}

// Our domain object for business logic
type WeatherService struct {
	Source Repo
}

// CurrentIn handles GET requests for finding current weather conditions at a latitude and longitude
func (w *WeatherService) CurrentIn(ctx context.Context, lat float32, lon float32) (*Weather, error) {
	cw, err := w.Source.GetByCoords(ctx, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("getting current weather by coordinates: %w", err)
	}
	// NOTE: this assumes we're using Imperial units, and is relative
	temp := TempUnknown
	switch {
	case cw.Temperature < 40.0:
		temp = TempCold
	case cw.Temperature < 80.0:
		temp = TempMod
	case cw.Temperature > 80.0:
		temp = TempHot
	}

	s := &Weather{
		Coords: types.Coords{
			Latitude:  cw.Coords.Latitude,
			Longitude: cw.Coords.Longitude,
		},
		States:      cw.States,
		Temperature: temp,
	}
	return s, nil
}

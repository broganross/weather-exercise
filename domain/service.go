package domain

import (
	"context"
	"fmt"
)

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
		Coords: Coords{
			Latitude:  cw.Coords.Latitude,
			Longitude: cw.Coords.Longitude,
		},
		States:      cw.States,
		Temperature: temp,
	}
	return s, nil
}

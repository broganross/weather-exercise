package domain

import (
	"context"
	"fmt"

	"github.com/broganross/weather-exercise/internal/repo"
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

// Interface for where we're getting actual weather data from
type WeatherSource interface {
	GetByCoords(ctx context.Context, latitude float32, longitude float32) (*repo.Weather, error)
}

// Our domain object for business logic
type WeatherService struct {
	Source WeatherSource
}

// CurrentIn handles GET requests for finding current weather conditions at a latitude and longitude
func (w *WeatherService) CurrentIn(ctx context.Context, lat float32, lon float32) (*Weather, error) {
	cw, err := w.Source.GetByCoords(ctx, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("getting current weather by coordinates: %w", err)
	}
	states := make([]string, len(cw.States))
	for i := 0; i < len(cw.States); i++ {
		states[i] = cw.States[i].Name
	}
	// NOTE: this assumes we're using Imperial units, and is relative
	temp := TempUnknown
	switch {
	case cw.Temperature.Temp < 40.0:
		temp = TempCold
	case cw.Temperature.Temp < 80.0:
		temp = TempMod
	case cw.Temperature.Temp > 80.0:
		temp = TempHot
	}

	s := &Weather{
		Coords: types.Coords{
			Latitude:  cw.Coords.Latitude,
			Longitude: cw.Coords.Longitude,
		},
		States:      states,
		Temperature: temp,
	}
	return s, nil
}

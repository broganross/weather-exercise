package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/broganross/weather-exercise/internal/domain"
	"github.com/broganross/weather-exercise/internal/types"
)

type WeatherState struct {
	ID          int
	Name        string
	Description string
}

type Temperature struct {
	Temp        float32
	FeelsLike   float32
	Min         float32
	Max         float32
	Pressure    int
	Humidity    int
	SeaLevel    int
	GroundLevel int
}

type Weather struct {
	Coords      types.Coords
	States      []WeatherState
	Temperature Temperature
}

type OpenWeather struct {
	BaseURL string
	Client  *http.Client
	APIid   string
	Timeout time.Duration
}

// GetByCoords retrieves current weather data for a set of coordinates
func (ow *OpenWeather) GetByCoords(ctx context.Context, lat float32, lon float32) (*domain.RepoWeather, error) {
	u := fmt.Sprintf("%s/weather", ow.BaseURL)
	ctx, cancel := context.WithTimeout(ctx, ow.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating open weather request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("lat", fmt.Sprintf("%02f", lat))
	q.Add("lon", fmt.Sprintf("%02f", lon))
	q.Add("appid", ow.APIid)
	req.URL.RawQuery = q.Encode()

	resp, err := ow.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var body string
		if b, err := io.ReadAll(resp.Body); err == nil {
			body = string(b)
		}
		err := fmt.Errorf("current weather by coordinates (%s): %s", http.StatusText(resp.StatusCode), body)
		return nil, err
	}
	item := currentWeatherResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("decoding current weather response body: %w", err)
	}
	states := make([]string, len(item.Weather))
	for index, w := range item.Weather {
		states[index] = w.Main
	}
	w := &domain.RepoWeather{
		Coords: types.Coords{
			Latitude:  item.Coord.Lat,
			Longitude: item.Coord.Lon,
		},
		States:      states,
		Temperature: item.Main.Temp,
	}
	return w, nil
}

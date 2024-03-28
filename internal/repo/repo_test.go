package repo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/broganross/weather-exercise/internal/repo"
	"github.com/broganross/weather-exercise/internal/types"
)

func TestOpenWeather_GetByCoords(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{
				"coord": {
				  "lat": 10.1,
				  "lon": 22.2
				},
				"weather": [
				  {
					"id": 501,
					"main": "Rain",
					"description": "moderate rain",
					"icon": "10d"
				  }
				],
				"base": "stations",
				"main": {
				  "temp": 298.48,
				  "feels_like": 298.74,
				  "temp_min": 297.56,
				  "temp_max": 300.05,
				  "pressure": 1015,
				  "humidity": 64,
				  "sea_level": 1015,
				  "grnd_level": 933
				},
				"visibility": 10000,
				"wind": {
				  "speed": 0.62,
				  "deg": 349,
				  "gust": 1.18
				},
				"rain": {
				  "1h": 3.16
				},
				"clouds": {
				  "all": 100
				},
				"dt": 1661870592,
				"sys": {
				  "type": 2,
				  "id": 2075663,
				  "country": "IT",
				  "sunrise": 1661834187,
				  "sunset": 1661882248
				},
				"timezone": 7200,
				"id": 3163858,
				"name": "Zocca",
				"cod": 200
			  }`))
		}))
	defer server.Close()
	want := &repo.Weather{
		Coords: types.Coords{
			Latitude:  10.1,
			Longitude: 22.2,
		},
		States: []repo.WeatherState{{
			ID:          501,
			Name:        "Rain",
			Description: "moderate rain",
		}},
		Temperature: repo.Temperature{
			Temp:        298.48,
			FeelsLike:   298.74,
			Min:         297.56,
			Max:         300.05,
			Pressure:    1015,
			Humidity:    64,
			SeaLevel:    1015,
			GroundLevel: 933,
		},
	}
	ow := repo.OpenWeather{
		BaseURL: server.URL,
		Client:  http.DefaultClient,
		APIid:   "API",
		Timeout: 5 * time.Second,
	}
	ctx := context.Background()
	got, err := ow.GetByCoords(ctx, 10.1, 22.2)
	if err != nil {
		t.Errorf("got unexpected error: '%v'", err)
		return
	}
	// Shouldn't actually use DeepEqual
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected '%v' got '%v'", want, got)
	}
}

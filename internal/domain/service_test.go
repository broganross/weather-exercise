package domain_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/broganross/weather-exercise/internal/domain"
	"github.com/broganross/weather-exercise/internal/repo"
	"github.com/broganross/weather-exercise/internal/types"
)

var errNotFound = errors.New("response not found")

var (
	rainState = repo.WeatherState{
		ID:          1,
		Name:        "rain",
		Description: "wetness falls from the sky",
	}
	hailState = repo.WeatherState{
		ID:          2,
		Name:        "hail",
		Description: "hard wetness falls from the sky",
	}
)

type mockWeatherRepoResponse struct {
	err  error
	resp *repo.Weather
}

type mockWeatherRepo struct {
	responses map[string]mockWeatherRepoResponse
}

func (mwr *mockWeatherRepo) GetByCoords(ctx context.Context, lat float32, lon float32) (*repo.Weather, error) {
	s := fmt.Sprintf("%.04f:%.04f", lat, lon)
	i, ok := mwr.responses[s]
	if !ok {
		return nil, errNotFound
	}
	return i.resp, i.err
}

func TestWeatherService_CurrentIn(t *testing.T) {
	tests := []struct {
		name string
		lat  float32
		lon  float32
		ans  *domain.Weather
		err  error
		repo mockWeatherRepo
	}{
		{
			"happy-path",
			10.1,
			32.1,
			&domain.Weather{
				Coords: types.Coords{
					Latitude:  10.1,
					Longitude: 32.1,
				},
				States:      []string{"rain", "hail"},
				Temperature: domain.TempCold,
			},
			nil,
			mockWeatherRepo{
				responses: map[string]mockWeatherRepoResponse{
					"10.1000:32.1000": {
						err: nil,
						resp: &repo.Weather{
							Coords: types.Coords{
								Latitude:  10.1,
								Longitude: 32.1,
							},
							States: []repo.WeatherState{
								rainState,
								hailState,
							},
							Temperature: repo.Temperature{
								Temp: 39.99,
							},
						},
					},
				},
			},
		},
		{
			"source-error",
			1.1,
			2.2,
			nil,
			errNotFound,
			mockWeatherRepo{
				responses: make(map[string]mockWeatherRepoResponse),
			},
		},
		// TODO: Extend tests
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serv := domain.WeatherService{
				Source: &test.repo,
			}
			got, err := serv.CurrentIn(context.Background(), test.lat, test.lon)
			if !errors.Is(err, test.err) {
				t.Errorf("expected '%v' got '%v'", test.err, err)
				return
			}
			// best not to use reflect, but it a fine placeholder for this exercise
			if !reflect.DeepEqual(test.ans, got) {
				t.Errorf("expected '%v' got '%v'", test.ans, *got)
			}
		})
	}
}

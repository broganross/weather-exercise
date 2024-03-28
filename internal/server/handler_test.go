package server_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broganross/weather-exercise/internal/domain"
	"github.com/broganross/weather-exercise/internal/server"
	"github.com/broganross/weather-exercise/internal/types"
)

var errResponseNotFound = errors.New("response not found")

type mockWeatherDomainResponse struct {
	weather domain.Weather
	err     error
}

type mockWeatherDomain struct {
	responses map[string]mockWeatherDomainResponse
}

func (mwd *mockWeatherDomain) CurrentIn(ctx context.Context, lat float32, lon float32) (*domain.Weather, error) {
	key := fmt.Sprintf("%.02f:%.02f", lat, lon)
	resp, ok := mwd.responses[key]
	if !ok {
		return nil, errResponseNotFound
	}
	return &resp.weather, resp.err
}

func TestWeatherSource_GetCurrentIn(t *testing.T) {
	handler := server.Handlers{
		Domain: &mockWeatherDomain{
			responses: map[string]mockWeatherDomainResponse{
				"1.20:2.30": {
					weather: domain.Weather{
						Coords: types.Coords{
							Latitude:  1.2,
							Longitude: 2.3,
						},
						States:      []string{"rain"},
						Temperature: domain.TempCold,
					},
				},
			},
		},
	}
	tests := []struct {
		name string
		path string
		code int
		body []byte
	}{
		{
			"happy-path",
			"?latitude=1.2&longitude=2.3",
			http.StatusOK,
			[]byte(`{"latitude":1.200000,"longitude":2.300000,"temperature":"cold","condition":"rain"}`),
		},
		{
			"missing-latitude",
			"?longitude=2.3",
			http.StatusBadRequest,
			[]byte(`{"error":"missing query parameter: latitude", "status":400}`),
		},
		{
			"missing-longitude",
			"?latitude=2.3",
			http.StatusBadRequest,
			[]byte(`{"error":"missing query parameter: longitude", "status":400}`),
		},
		{
			"domain-failure",
			"?latitude=1.22&longitude=2.30",
			http.StatusInternalServerError,
			[]byte(`{"error":"retrieving current weather: no response found","status":500}`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost/"+test.path, nil)
			w := httptest.NewRecorder()
			handler.GetCurrentByCoords(w, req)
			resp := w.Result()
			if resp.StatusCode != test.code {
				t.Errorf("expected code '%v' got '%v'", test.code, resp.StatusCode)
			}
			if h := resp.Header.Get("Content-Type"); h != "application/json" {
				t.Errorf("expected Content-Type header 'application/json' got '%v'", h)
			}
			body, _ := io.ReadAll(resp.Body)
			if bytes.Equal(test.body, body) {
				t.Errorf("expected body '%v' got '%v'", test.body, string(body))
			}
		})
	}
}

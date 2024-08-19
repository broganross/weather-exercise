package server_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/broganross/weather-exercise/domain"
	"github.com/broganross/weather-exercise/server"
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
						Coords: domain.Coords{
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
			[]byte("{\"id\":\"urn:weather:current:id\",\"type\":\"urn:weather:current\",\"attributes\":{\"latitude\":1.200000,\"longitude\":2.300000,\"temperature\":\"cold\",\"condition\":\"rain\"}}\n"),
		},
		{
			"missing-parameters",
			"",
			http.StatusBadRequest,
			[]byte("{\"errors\":[{\"error\":\"missing query parameter: latitude\",\"message\":\"required query parameters\"},{\"error\":\"missing query parameter: longitude\",\"message\":\"required query parameters\"}],\"status\":400}\n"),
		},
		{
			"domain-failure",
			"?latitude=1.22&longitude=2.30",
			http.StatusInternalServerError,
			[]byte("{\"errors\":[{\"error\":\"retrieving current weather: response not found\"}],\"status\":500}\n"),
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
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("reading body: %v", err)
			}
			s := string(body)
			_ = s
			e := string(test.body)
			_ = e
			if !strings.EqualFold(string(test.body), s) {
				// if !bytes.Equal(test.body, body) {
				t.Errorf("expected body '%v' got '%v'", string(test.body), s)
			}
		})
	}
}

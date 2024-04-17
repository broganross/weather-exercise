package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/broganross/weather-exercise/internal/config"
	"github.com/broganross/weather-exercise/internal/domain"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var (
	ErrMissingParam = errors.New("missing query parameter")
	ErrInvalidFloat = errors.New("invalid float")
)

// SetupRoutes constructs the router, adding middleware, routes, handlers, etc
func SetupRoutes(h *Handlers, r *mux.Router, c *config.Config) {
	r.Use(LogContextMiddleware)
	am := Auth{
		BaseURL: c.AuthService.URL,
	}
	r.Use(am.Middleware)
	r.HandleFunc("/", h.GetCurrentByCoords).Methods(http.MethodGet)
}

// Our handlers for whatever routes we need
type Handlers struct {
	Domain domain.Service
}

func (h *Handlers) GetCurrentByCoords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	q := r.URL.Query()
	// get the query parameters
	var lat float64
	var lon float64
	latString := q.Get("latitude")
	if latString == "" {
		encodeError(ctx, w, http.StatusBadRequest, fmt.Errorf("%w: latitude", ErrMissingParam), "")
		return
	}
	lat, err := strconv.ParseFloat(latString, 32)
	if err != nil {
		encodeError(ctx, w, http.StatusBadRequest, fmt.Errorf("%w: latitude", ErrInvalidFloat), "")
		return
	}
	lonString := q.Get("longitude")
	if lonString == "" {
		encodeError(ctx, w, http.StatusBadRequest, fmt.Errorf("%w: longitude", ErrMissingParam), "")
		return
	}
	lon, err = strconv.ParseFloat(lonString, 32)
	if err != nil {
		encodeError(ctx, w, http.StatusBadRequest, fmt.Errorf("%w: longitude", ErrInvalidFloat), "")
		return
	}

	// business logic
	weather, err := h.Domain.CurrentIn(ctx, float32(lat), float32(lon))
	if err != nil {
		encodeError(
			ctx,
			w,
			http.StatusInternalServerError,
			fmt.Errorf("retrieving current weather: %w", err),
			"",
		)
		return
	}
	// remap structure to API
	cond := strings.Join(weather.States, ", ")
	resp := getCurrentByCoordsResponse{
		Temperature: string(weather.Temperature),
		Condition:   cond,
		Latitude:    preciseFloat32(lat),
		Longitude:   preciseFloat32(lon),
	}
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		encodeError(
			ctx,
			w,
			http.StatusInternalServerError,
			fmt.Errorf("encoding current weather response: %w", err),
			"",
		)
		return
	}
}

// Creates and writes an error
func encodeError(ctx context.Context, w http.ResponseWriter, statusCode int, err error, message string) {
	l := log.Ctx(ctx)
	item := errorResponse{
		Error:   err.Error(),
		Message: message,
		Status:  statusCode,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&item); err != nil {
		b := []byte(fmt.Sprintf(`{"status":500,"error":"error encoding error response: %s}`, err))
		if _, err := w.Write(b); err != nil {
			l.Error().Err(err).Msg("writing default error to response writer")
			return
		}
	}
	l.Err(err).Int("status_code", statusCode).Send()
}

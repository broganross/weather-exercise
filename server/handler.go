package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/broganross/weather-exercise/config"
	"github.com/broganross/weather-exercise/domain"
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
	errs := []error{}
	if latString := q.Get("latitude"); latString != "" {
		l, err := strconv.ParseFloat(latString, 32)
		if err != nil {
			errs = append(errs, fmt.Errorf("%w: latitude", ErrInvalidFloat))
		}
		lat = l
	} else {
		errs = append(errs, fmt.Errorf("%w: latitude", ErrMissingParam))
	}
	if lonString := q.Get("longitude"); lonString != "" {
		l, err := strconv.ParseFloat(lonString, 32)
		if err != nil {
			errs = append(errs, fmt.Errorf("%w: longitude", ErrInvalidFloat))
		}
		lon = l
	} else {
		errs = append(errs, fmt.Errorf("%w: longitude", ErrMissingParam))
	}
	if len(errs) > 0 {
		encodeError(ctx, w, http.StatusBadRequest, errs, "required query parameters")
		return
	}

	// business logic
	weather, err := h.Domain.CurrentIn(ctx, float32(lat), float32(lon))
	if err != nil {
		encodeError(
			ctx,
			w,
			http.StatusInternalServerError,
			[]error{fmt.Errorf("retrieving current weather: %w", err)},
			"",
		)
		return
	}
	// remap structure to API
	cond := strings.Join(weather.States, ", ")
	resp := getCurrentByCoordsResponse{
		ID:   "urn:weather:current:id",
		Type: "urn:weather:current",
		Attribtues: currentAttributes{
			Temperature: string(weather.Temperature),
			Condition:   cond,
			Latitude:    preciseFloat32(lat),
			Longitude:   preciseFloat32(lon),
		},
	}
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		encodeError(
			ctx,
			w,
			http.StatusInternalServerError,
			[]error{fmt.Errorf("encoding current weather response: %w", err)},
			"",
		)
		return
	}
}

// Creates and writes an error
func encodeError(ctx context.Context, w http.ResponseWriter, statusCode int, errs []error, message string) {
	l := log.Ctx(ctx)
	resp := errorResponse{Status: statusCode}
	event := l.Error()
	for _, e := range errs {
		item := errorItem{
			Error:   e.Error(),
			Message: message,
		}
		resp.Errors = append(resp.Errors, item)
		event.Err(e)
	}
	event.Int("status_code", statusCode).Send()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		b := []byte(fmt.Sprintf(`{"status":500,"error":"error encoding error response: %s}`, err))
		if _, err := w.Write(b); err != nil {
			l.Error().Err(err).Msg("writing default error to response writer")
			return
		}
	}
}

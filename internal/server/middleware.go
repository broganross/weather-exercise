package server

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogContextMiddleware injects a logger into the context and adds a request id.
// TODO: extend to log outgoing request statuses
func LogContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		l := log.With().Logger()
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request-id", requestID)
		})
		r = r.WithContext(l.WithContext(r.Context()))
		l.Debug().
			Str("method", r.Method).
			Stringer("URL", r.URL).
			Str("request-id", requestID).
			Msg("incoming request")
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware an example auth external service
type Auth struct {
	BaseURL string
}

func (am *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticated := true
		authorized := true
		// pull an auth token from the header
		// check it's validity in some repository
		// check that the user can access the resource
		if !authenticated || !authorized {
			var err error
			if !authenticated {
				err = errors.New("user not authenticated")
			} else if !authorized {
				err = errors.New("user not authorized")
			}
			EncodeError(
				r.Context(),
				w,
				http.StatusUnauthorized,
				err,
				"",
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Debug().Msg("this is where we would be checking authentication and authorization")
// 		next.ServeHTTP(w, r)
// 	})
// }

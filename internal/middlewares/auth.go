package middlewares

import (
	"net/http"
	"wave/internal/service"
)

// easyjson:json
type authConfig struct {
	Adress string `json:"adress"`
}

// Auth - authorisation check middleware
func Auth(sv service.IService) IMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")

			if err != nil || cookie.Value == "" {
				sv.Logger().Info("Auth check failed") //TODO:: usless log. modify or remove
				w.WriteHeader(http.StatusUnauthorized)
				//TODO:: metrics
				return
			}

			// TODO:: check auth

			next(w, r)
		}
	}
}

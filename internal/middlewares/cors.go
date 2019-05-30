package middlewares

import (
	"fmt"
	"net/http"
	"wave/internal/service"

	"github.com/rs/cors"
)

// easyjson:json
type corsConfig struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Methods     []string `json:"methods"`
	Credentials bool     `json:"credentials"`
	OptionsPass bool     `json:"optionspass"`
}

// Cors - cors middleware
func Cors(sv service.IService) IMiddleware {
	config := &corsConfig{}
	if err := sv.Config(config, "cors"); err != nil {
		service.Panic(fmt.Errorf(`Unexpected error during cors initialisation`), err)
	}
	cors := cors.New(cors.Options{
		OptionsPassthrough: config.OptionsPass,
		AllowCredentials:   config.Credentials,
		AllowedOrigins:     config.Origins,
		AllowedHeaders:     config.Headers,
		AllowedMethods:     config.Headers,
	})
	return func(next http.HandlerFunc) http.HandlerFunc {
		handler := cors.Handler(http.HandlerFunc(next))
		return handler.ServeHTTP
	}
}

package middlewares

import "net/http"

//go:generate easyjson -output_filename jsons.go .

// IMiddleware - middleware inteface
type IMiddleware func(http.HandlerFunc) http.HandlerFunc

// Pipe - cretae a middleware pipe
func Pipe(handler http.HandlerFunc, middlewares ...IMiddleware) http.HandlerFunc {
	for _, mw := range middlewares {
		handler = mw(handler)
	}
	return handler
}

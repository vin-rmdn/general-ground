package middleware

import "net/http"

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error
type Middleware func(next HandlerWithError) HandlerWithError

func Wrap(next HandlerWithError, middlewares ...Middleware) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, middleware := range middlewares {
			next = middleware(next)
		}

		_ = next(w, r)
	})
}
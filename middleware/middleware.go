package middleware

import "net/http"

type Middleware = func(http.Handler) http.Handler

func NoOp() Middleware {
	return func(next http.Handler) http.Handler {
		return next
	}
}

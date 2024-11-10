package middleware

import "net/http"

// Middleware определяет тип функции middleware
type Middleware func(http.Handler) http.Handler

// ChainMiddleware создаёт цепочку middleware
func ChainMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

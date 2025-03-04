package httpserver

import (
	"net/http"
	"strings"
)

// App is a HTTP server
type App struct {
	mux         *http.ServeMux
	middlewares []MiddlewareFunc
	server      *http.Server
}

// New creates a new App
func New() *App {
	mux := http.NewServeMux()
	svr := &http.Server{}

	return &App{mux: mux, server: svr}
}

// Get adds a GET route to the server
//   - path: the route path
//   - handler: the route handler func(w http.ResponseWriter, r *http.Request)
func (a *App) Get(path string, handler http.HandlerFunc) {
	pattern := strings.Join([]string{http.MethodGet, path}, " ")
	a.mux.HandleFunc(pattern, handler)
}

// Post adds a POST route to the server
//   - path: the route path
//   - handler: the route handler func(w http.ResponseWriter, r *http.Request)
func (a *App) Post(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	pattern := strings.Join([]string{http.MethodPost, path}, " ")
	a.mux.HandleFunc(pattern, handler)
}

// Put adds a PUT route to the server
//   - path: the route path
//   - handler: the route handler func(w http.ResponseWriter, r *http.Request)
func (a *App) Put(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	pattern := strings.Join([]string{http.MethodPut, path}, " ")
	a.mux.HandleFunc(pattern, handler)
}

// Delete adds a DELETE route to the server
//   - path: the route path
//   - handler: the route handler func(w http.ResponseWriter, r *http.Request)
func (a *App) Delete(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	pattern := strings.Join([]string{http.MethodDelete, path}, " ")
	a.mux.HandleFunc(pattern, handler)
}

// Patch adds a PATCH route to the server
//   - path: the route path
//   - handler: the route handler func(w http.ResponseWriter, r *http.Request)
func (a *App) Patch(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	pattern := strings.Join([]string{http.MethodPatch, path}, " ")
	a.mux.HandleFunc(pattern, handler)
}

// Listen starts the server with given port
//   - port: the port to listen
func (a *App) Listen(port string) error {
	a.server.Addr = ":" + port
	a.server.Handler = a.mux

	// Apply middlewares
	for _, f := range a.middlewares {
		a.server.Handler = f(a.server.Handler)
	}

	return a.server.ListenAndServe()
}

// Use adds middlewares to the server.
//   - MiddlewareFunc is a HTTP middleware function.
//   - type MiddlewareFunc func(http.Handler) http.Handler
func (a *App) Use(middlewares ...MiddlewareFunc) {
	for _, m := range middlewares {
		a.middlewares = append(a.middlewares, m)
	}
}

// MiddlewareFunc is a HTTP middleware function
type MiddlewareFunc func(http.Handler) http.Handler

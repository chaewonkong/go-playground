package httpserver

import "net/http"

// Server is a HTTP server
type Server struct {
	server *http.Server
}

// New creates a new server
func New(mux *http.ServeMux, svr *http.Server) *Server {
	if svr == nil {
		svr = &http.Server{
			Addr: ":8080",
		}
	}
	svr.Handler = mux
	return &Server{
		server: svr,
	}
}

// Use adds middlewares to the server
func (s *Server) Use(middlewares ...MiddlewareFunc) {
	h := s.server.Handler
	for _, m := range middlewares {
		h = m(h)
	}
	s.server.Handler = h
}

// Run starts the server
func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

// MiddlewareFunc is a HTTP middleware function
type MiddlewareFunc func(http.Handler) http.Handler

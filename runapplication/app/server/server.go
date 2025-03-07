package server

import (
	"net/http"
	"runapplication/app"
)

// Config 구조체
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`
	Database struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
	} `mapstructure:"database"`
}

// Server 서버 구조체
type Server struct {
	svr    *http.Server
	config *Config
}

// New 서버 생성자
func New(cfg *Config) app.App {
	svr := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: newMux(),
	}

	return &Server{
		svr:    svr,
		config: cfg,
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /greet", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}

	})
	return mux
}

// Run 서버 실행
func (s *Server) Run() error {
	return s.svr.ListenAndServe()
}

// Stop 서버 종료
func (s *Server) Stop() error {
	return s.svr.Close()
}

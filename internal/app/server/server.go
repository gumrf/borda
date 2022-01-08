package server

import (
	"context"
	"fmt"
	"net/http"
	"borda/internal/app/setup"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, config setup.HTTPConfig) *Server {
	serverAddr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	return &Server{
		httpServer: &http.Server{Addr: serverAddr, Handler: handler},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

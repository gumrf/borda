package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger *zap.Logger
}

func NewServer(handler http.Handler, logger *zap.Logger) *Server {
	return &Server{
		httpServer: &http.Server{Addr: ":8080", Handler: handler},
		logger:     logger,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

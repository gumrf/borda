package server

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, host string, port string) *Server {
	serverPort := fmt.Sprintf("%s:%s", host, port)
	return &Server{
		httpServer: &http.Server{Addr: serverPort, Handler: handler},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

package core

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func (srv *Server) Run(port string, handler http.Handler) error {
	srv.server = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.server.ListenAndServe()
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.server.Shutdown(ctx)
}

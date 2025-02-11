// Package serverapp defines server application struct that accepts connections.
package serverapp

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strconv"
)

var errServerStopped = errors.New("server has stopped")

// Server listens on port for new http connections and passes them to router.
type Server struct {
	log    *slog.Logger
	router *gin.Engine
	port   int
}

// New returns new server instance.
func New(
	log *slog.Logger,
	port int,
	router *gin.Engine,
) *Server {
	return &Server{
		log:    log,
		port:   port,
		router: router,
	}
}

// Serve starts http server and returns error when server is stopped.
func (s *Server) Serve() error {
	const op = "serverApp.Serve"

	log := s.log.With(slog.String("op", op))

	addr := fmt.Sprintf(":%s", strconv.Itoa(s.port))

	log.Info(fmt.Sprintf("listening on %s", addr))

	err := s.router.Run(addr)

	if err == nil {
		return nil
	}

	log.Error("server has stopped", slog.Any("error", err.Error()))

	return fmt.Errorf("%w: %s", errServerStopped, err.Error())
}

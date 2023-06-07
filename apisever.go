package contester_go

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
}

func (s *Server) StartServer(addr string, handler http.Handler) error {
	s.server = &http.Server{
		Addr:           addr,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   1 * time.Minute,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
	}

	log.Info("server started")

	return s.server.ListenAndServe()
}

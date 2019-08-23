package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/rancher/webhook-receiver/pkg/apis"
	"github.com/rancher/webhook-receiver/pkg/options"
)

type Server struct {
	port       int
	configPath string
}

func New(port int, configPath string) *Server {
	return &Server{
		port:       port,
		configPath: configPath,
	}
}

func (s *Server) Run() error {
	options.Init(s.configPath)
	apis.RegisterAPIs()

	log.Infof("server running, listening at :%d\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

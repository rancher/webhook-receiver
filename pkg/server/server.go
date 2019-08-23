package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rancher/receiver/pkg/apis"
	"github.com/rancher/receiver/pkg/options"
)

type Server struct {
	port       int
	configPath string
}

func New(port int, configPath string) (*Server, error) {
	s := &Server{
		port:       port,
		configPath: configPath,
	}

	return s, nil
}

func (s *Server) Run() {
	options.Init(s.configPath)
	apis.RegisterAPIs()

	log.Printf("server running, listening at :%d\n", s.port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

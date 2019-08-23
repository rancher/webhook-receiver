package main

import (
	"flag"

	"github.com/rancher/receiver/pkg/server"

	log "github.com/sirupsen/logrus"
)

var (
	port   = flag.Int("port", 9094, "listen port")
	config = flag.String("config", "/etc/receiver/config.yaml", "config path")
)

func main() {
	flag.Parse()
	if err := server.New(*port, *config).Run(); err != nil {
		log.Fatalf("server run fatal:%v", err)
	}
}

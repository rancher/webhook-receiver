package main

import (
	"flag"
	"log"

	"github.com/rancher/receiver/pkg/server"
)

var (
	port = flag.Int("port", 9094, "listen port")
	config = flag.String("config", "/etc/receiver/config.yaml", "config path")
)

func main() {
	flag.Parse()
	s, err := server.New(*port, *config)
	if err != nil {
		log.Printf("new server error:%v", err)
	}
	s.Run()
}

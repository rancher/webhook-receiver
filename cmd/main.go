package main

import (
	"flag"
	"log"

	"github.com/rancher/webhook-receiver/pkg/server"
)

var (
	port   = flag.Int("port", 9094, "listen port")
	config = flag.String("config", "/etc/webhook-receiver/config.yaml", "config path")
)

func main() {
	flag.Parse()
	if err := server.New(*port, *config).Run(); err != nil {
		log.Fatalf("server run fatal:%v", err)
	}
}

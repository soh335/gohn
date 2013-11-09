package main

import (
	"flag"
	"github.com/soh335/gohn"
)

var (
	host        = flag.String("host", "127.0.0.1", "host")
	port        = flag.String("port", "5555", "port")
	rpcHost     = flag.String("rpcHost", "127.0.0.1", "rpcHost")
	rpcPort     = flag.String("rpcPort", "5556", "rpcPort")
	rpcParallel = flag.Int("rpcParallel", 10, "rpcParallel")
	config      = flag.String("config", "", "config json")
)

func main() {
	flag.Parse()

	configLoader := gohn.NewConfigLoader(*config)
	go configLoader.Start(*rpcHost, *rpcPort, *rpcParallel)
	go gohn.StartHttpServer(*host, *port, configLoader)
	select {}
}

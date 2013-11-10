package main

import (
	"flag"
	"log"
	"os"
)

var (
	dataDir     = flag.String("datadir", "./data", "datadir")
	host        = flag.String("host", "127.0.0.1", "host")
	port        = flag.String("port", "5555", "port")
	rpcHost     = flag.String("rpcHost", "127.0.0.1", "rpcHost")
	rpcPort     = flag.String("rpcPort", "5556", "rpcPort")
	rpcParallel = flag.Int("rpcParallel", 10, "rpcParallel")
	config      = flag.String("config", "", "config json")
	playCmd     = flag.String("playCmd", "afplay", "play cmd")
)

func main() {
	flag.Parse()

	err := os.MkdirAll(*dataDir, 0777)
	if err != nil {
		log.Fatal("crete dir err", *dataDir, err)
	}

	go StartRpcServer(*rpcHost, *rpcPort, *dataDir, *playCmd)

	configLoader := NewConfigLoader(*config)
	go configLoader.Start(*rpcHost, *rpcPort, *rpcParallel)
	go StartHttpServer(*host, *port, configLoader)

	select {}
}

package main

import (
	"flag"
	"github.com/soh335/gohn"
	"log"
	"os"
)

var (
	host    = flag.String("host", "127.0.0.1", "host")
	port    = flag.String("port", "5556", "port")
	dataDir = flag.String("datadir", "./data", "datadir")
)

func main() {
	flag.Parse()

	err := os.MkdirAll(*dataDir, 0777)
	if err != nil {
		log.Fatal("crete dir err", *dataDir, err)
	}

	go gohn.StartRpcServer(*host, *port, *dataDir)
	select {}
}

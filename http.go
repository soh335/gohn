package gohn

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

type Manager int

type TagM4A struct {
	M4A M4A
	Tag string
}

type Config map[string][]Source

type ConfigLoader struct {
	configPath   string
	requestChan  chan string
	responseChan chan *M4A
	client       *rpc.Client
}

func NewConfigLoader(configPath string) *ConfigLoader {
	configLoader := &ConfigLoader{}
	configLoader.configPath = configPath
	configLoader.requestChan = make(chan string)
	configLoader.responseChan = make(chan *M4A)

	return configLoader
}

func StartHttpServer(host string, port string, configLoader *ConfigLoader) {

	http.HandleFunc("/play/", func(w http.ResponseWriter, req *http.Request) {
		playHandle(w, req, configLoader)
	})

	addr := net.JoinHostPort(host, port)
	log.Println("starting httpd", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func playHandle(w http.ResponseWriter, req *http.Request, configLoader *ConfigLoader) {

	tag := req.URL.Path[6:]
	log.Println(tag)
	configLoader.requestChan <- tag

	m4a := <-configLoader.responseChan
	if m4a != nil {
		var res PlayResponse
		err := configLoader.client.Call("Executor.Play", *m4a, &res)
		if err != nil {
			log.Println("play call", err)
		}
		return
	} else {
		log.Println("not found")
		return
	}
}

func (c *ConfigLoader) Start(host string, port string, rpcParallel int) {
	addr := net.JoinHostPort(host, port)
	var client *rpc.Client
	var err error
	for {
		client, err = jsonrpc.Dial("tcp", addr)
		if err != nil {
			log.Println("rpc connect", err)
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
	log.Println("connect to rcp server", addr)

	c.client = client

	tagM4AChan := make(chan TagM4A)
	convertedTagM4Achan := make(chan TagM4A)
	semaphore := make(chan bool, rpcParallel)

	go func() {
		m4AMapList := make(map[string][]M4A)
		for {
			select {
			case tagM4A := <-tagM4AChan:
				log.Println(tagM4A)
				if _, ok := m4AMapList[tagM4A.Tag]; ok == false {
					m4AMapList[tagM4A.Tag] = make([]M4A, 0)
				}
				m4AMapList[tagM4A.Tag] = append(m4AMapList[tagM4A.Tag], tagM4A.M4A)
				log.Println(m4AMapList)
			case tag := <-c.requestChan:
				m4AList, ok := m4AMapList[tag]
				if !(ok == true && len(m4AList) > 0) {
					c.responseChan <- nil
					continue
				}
				m4a := m4AList[rand.Intn(len(m4AList))]
				c.responseChan <- &m4a
			}
		}
	}()

	go func() {
		config := openConfig(c.configPath)
		for tag, sourceList := range *config {
			for _, source := range sourceList {
				semaphore <- true
				go c.startConvert(tag, source, convertedTagM4Achan)
			}
		}
	}()

	go func() {
		for {
			<-semaphore
			tagM4AChan <- <-convertedTagM4Achan
		}
	}()
}

func (c *ConfigLoader) startConvert(tag string, source Source, convertedTagM4Achan chan TagM4A) {
	// executor に投げる
	log.Println(tag, source)
	var m4a M4A
	err := c.client.Call("Executor.Convert", source, &m4a)
	if err != nil {
		log.Println("Executor.Convert err ", err)
	}

	convertedTagM4Achan <- TagM4A{m4a, tag}
}

func openConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		log.Fatal(err)
	}

	return &config
}

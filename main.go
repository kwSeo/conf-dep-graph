package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/kwseo/conf-dep-graph/pkg/topology"
)

func main() {
	var configFile string
	var serviceName string
	flag.StringVar(&configFile, "config-file", "", "configuration file")
	flag.StringVar(&serviceName, "service", "", "root service name to create a dependency graph")
	flag.Parse()

	cfg, err := topology.LoadConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	t := topology.New(cfg)
	for _, svc := range cfg.Services {
		if _, err := svc.LoadTargetFile(); err != nil {
			log.Fatalln(err)
		}
		t.AddService(svc)
	}

	var png []byte
	if serviceName == "" {
		png, err = t.GraphAsPNG()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		png, err = t.ServiceGraphAsPNG(serviceName)
		if err != nil {
			log.Fatalln(err)
		}
	}
	if err := ioutil.WriteFile("graph.png", png, os.ModePerm); err != nil {
		log.Fatalln(err)
	}
}

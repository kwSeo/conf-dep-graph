package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kwseo/conf-dep-graph/pkg/topology"
)

var logger log.Logger = log.With(log.NewLogfmtLogger(os.Stdout), "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

func main() {
	var configFile string
	var serviceName string
	flag.StringVar(&configFile, "config-file", "", "configuration file")
	flag.StringVar(&serviceName, "service", "", "root service name to create a dependency graph")
	flag.Parse()

	cfg, err := topology.LoadConfig(configFile)
	if err != nil {
		fatal(err)
	}

	t := topology.New(cfg)
	for _, svc := range cfg.Services {
		if _, err := svc.LoadTargetFile(); err != nil {
			fatal(err)
		}
		t.AddService(svc)
	}

	var png []byte
	if serviceName == "" {
		png, err = t.GraphAsPNG()
		if err != nil {
			fatal(err)
		}
	} else {
		png, err = t.ServiceGraphAsPNG(serviceName)
		if err != nil {
			fatal(err)
		}
	}
	if err := ioutil.WriteFile("graph.png", png, os.ModePerm); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	level.Error(logger).Log("err", err)
	os.Exit(1)
}

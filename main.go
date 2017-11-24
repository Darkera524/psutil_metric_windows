package main

import (
	"github.com/Darkera524/psutil_metric_windows/cron"
	"flag"
	"github.com/golang/glog"
	"github.com/Darkera524/psutil_metric_windows/g"
)

var cfg = flag.String("c", "cfg.example.json", "configuration file")
var version = flag.Bool("version", false, "show version")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	defer glog.Flush()

	flag.Parse()

	g.HandleVersion(*version)
	if memfile, _ := g.HandleMemProfile(*memprofile); memfile != nil {
		defer memfile.Close()
	}

	// global config
	g.ParseConfig(*cfg)
	g.InitRpcClients()

	cron.Collect()
}

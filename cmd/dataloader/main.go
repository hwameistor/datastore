package main

import (
	"flag"

	"github.com/hwameistor/datastore/pkg/dataloader"

	log "github.com/sirupsen/logrus"
)

func setupLogging() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	flag.Parse()
	setupLogging()

	stopCh := make(chan struct{})
	dataloader.NewManager().Run(stopCh)

}

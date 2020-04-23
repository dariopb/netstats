package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"

	restapi "github.com/dariopb/netstats/pkg/restapi"
)

func main() {
	var err error = nil
	var port int = 9090
	portstr, ok := os.LookupEnv("PORT")
	if ok {
		port, err = strconv.Atoi(portstr)
		if err != nil {
			panic(err)
		}
	}

	log.Infof("Starting netstats service on port %d\n", port)

	restapi.NewRestApi(port)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}

package main

import (
	"beryju.io/saml-test-sp/pkg/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	server.RunServer()
}

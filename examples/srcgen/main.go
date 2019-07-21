package main

import (
	"github.com/muka/go-bluetooth/gen"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("Generating src")
	gen.Parse("/home/l/git/kernel.org/bluetooth/bluez")
	gen.Generate("./tmp")

}

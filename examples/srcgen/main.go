package main

import (
	"github.com/muka/go-bluetooth/gen"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("Generating src")
	api := gen.Parse("/home/l/git/kernel.org/bluetooth/bluez")
	// api := gen.Parse("./test")
	gen.Generate(api, "./test/out")

}

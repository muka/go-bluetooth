package main

import (
	"os"

	"github.com/muka/go-bluetooth/gen"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("Generating src")

	var api []gen.ApiGroup
	if os.Getenv("DEBUG") != "" {
		api = gen.Parse("./test")
	} else {
		api = gen.Parse("/home/l/git/kernel.org/bluetooth/bluez")
	}

	gen.Generate(api, "./test/out")

}

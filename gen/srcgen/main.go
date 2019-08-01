package main

import (
	"os"

	"github.com/muka/go-bluetooth/gen"
	"github.com/muka/go-bluetooth/gen/generator"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("Generating src")

	var api []gen.ApiGroup
	if os.Getenv("DEBUG") == "1" {
		api = gen.Parse("./test")
		generator.Generate(api, "./src/gen")
	} else {
		api = gen.Parse("./src/bluez")
		generator.Generate(api, "./bluez")
	}

}

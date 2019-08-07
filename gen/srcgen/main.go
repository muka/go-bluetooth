package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/muka/go-bluetooth/gen"
	"github.com/muka/go-bluetooth/gen/generator"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("Generating src")

	filters := strings.Split(os.Getenv("FILTER"), ",")

	api, err := gen.Parse("./src/bluez/doc", filters)
	if err != nil {
		log.Fatalf("Parse failed: %s", err)
	}

	err = api.Serialize(fmt.Sprintf("./bluez-%s.json", api.Version))
	if err != nil {
		log.Fatalf("Failed to serialize JSON: %s", err)
	}

	err = generator.Generate(api, "./bluez", false)
	if err != nil {
		log.Fatalf("Generation failed: %s", err)
	}

}

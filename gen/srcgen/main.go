// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/muka/go-bluetooth/gen"
	"github.com/muka/go-bluetooth/gen/generator"
	log "github.com/sirupsen/logrus"
)

const ApiFile = "./bluez-5.50.json"

func main() {

	log.SetLevel(log.DebugLevel)

	if len(os.Args) > 1 && os.Args[1] == "full" {
		log.Info("Generating src")
		filters := strings.Split(os.Getenv("FILTER"), ",")
		err := Parse(filters)
		if err != nil {
			os.Exit(1)
		}
	}

	err := Generate(ApiFile)
	if err != nil {
		os.Exit(1)
	}

}

func Generate(filename string) error {

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Generation failed: %s", err)
		return err
	}

	api := gen.BluezAPI{}
	err = json.Unmarshal([]byte(file), &api)
	if err != nil {
		log.Fatalf("Generation failed: %s", err)
		return err
	}

	err = generator.Generate(api, "./bluez", false)
	if err != nil {
		log.Fatalf("Generation failed: %s", err)
		return err
	}

	return nil
}

func Parse(filters []string) error {

	api, err := gen.Parse("./src/bluez/doc", filters)
	if err != nil {
		log.Fatalf("Parse failed: %s", err)
		return err
	}

	err = api.Serialize(fmt.Sprintf("./bluez-%s.json", api.Version))
	if err != nil {
		log.Fatalf("Failed to serialize JSON: %s", err)
		return err
	}

	return nil
}

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

const docsDir = "./src/bluez/doc"

func main() {

	log.SetLevel(log.DebugLevel)

	bluezVersion, err := gen.GetGitVersion(docsDir)
	if err != nil {
		log.Fatal(err)
	}

	envBluezVersion := os.Getenv("BLUEZ_VERSION")
	if envBluezVersion != "" {
		bluezVersion = envBluezVersion
	}

	fmt.Printf("---\nAPI %s\n---\n", bluezVersion)

	apiFile := fmt.Sprintf("./bluez-%s.json", bluezVersion)

	if len(os.Args) > 1 && os.Args[1] == "full" {
		log.Info("Generating src")
		filters := strings.Split(os.Getenv("FILTER"), ",")
		err := Parse(filters)
		if err != nil {
			os.Exit(1)
		}
	}

	err = Generate(apiFile)
	if err != nil {
		log.Fatal(err)
	}

}

func Generate(filename string) error {

	fmt.Printf("Generating from %s\n", filename)

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

	api, err := gen.Parse(docsDir, filters)
	if err != nil {
		log.Fatalf("Parse failed: %s", err)
		return err
	}

	apiFile := fmt.Sprintf("./bluez-%s.json", api.Version)
	fmt.Printf("Creating from %s\n", apiFile)
	err = api.Serialize(apiFile)
	if err != nil {
		log.Fatalf("Failed to serialize JSON: %s", err)
		return err
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/muka/go-bluetooth/gen"
	"github.com/muka/go-bluetooth/gen/generator"
	"github.com/muka/go-bluetooth/gen/util"
	log "github.com/sirupsen/logrus"
)

const docsDir = "./src/bluez/doc"

func main() {

	logLevel := log.DebugLevel.String()
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)

	bluezVersion, err := util.GetGitVersion(docsDir)
	if err != nil {
		log.Fatal(err)
	}

	envBluezVersion := os.Getenv("BLUEZ_VERSION")
	if envBluezVersion != "" {
		bluezVersion = envBluezVersion
	}

	generateMode := "full"
	if len(os.Args) > 1 {
		generateMode = os.Args[1]
	}

	log.Infof("API %s", bluezVersion)

	apiFile := fmt.Sprintf("./bluez-%s.json", bluezVersion)

	if generateMode == "full" || generateMode == "parse" {
		log.Debug("Generating src")
		filters := strings.Split(os.Getenv("FILTER"), ",")
		err := Parse(filters, lvl == log.DebugLevel)
		if err != nil {
			os.Exit(1)
		}
	}

	if generateMode == "full" || generateMode == "generate" {
		err = Generate(apiFile)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func Parse(filters []string, debug bool) error {

	api, err := gen.Parse(docsDir, filters, debug)
	if err != nil {
		log.Fatalf("Parse failed: %s", err)
		return err
	}

	apiFile := fmt.Sprintf("./bluez-%s.json", api.Version)
	log.Infof("Creating from %s\n", apiFile)
	err = api.Serialize(apiFile)
	if err != nil {
		log.Fatalf("Failed to serialize JSON: %s", err)
		return err
	}

	return nil
}

func Generate(filename string) error {

	log.Infof("Generating from %s\n", filename)

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

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
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	flagOverwrite            = "overwrite"
	flagDebug                = "debug"
	flagGenerateModeFull     = "full"
	flagGenerateModeParse    = "parse"
	flagGenerateModeGenerate = "generate"
	paramFilter              = "filter"
)

const docsDir = "./src/bluez/doc"

func main() {

	parseLogLevel()
	bluezVersion := getBluezVersion()
	debug := hasFlag(flagDebug)

	apiFile := fmt.Sprintf("./bluez-%s.json", bluezVersion)

	if hasFlag(flagGenerateModeFull) || hasFlag(flagGenerateModeParse) {
		filters := parseFilters()
		err := Parse(filters, debug)
		if err != nil {
			os.Exit(1)
		}
	}

	if hasFlag(flagGenerateModeFull) || hasFlag(flagGenerateModeGenerate) {
		overwrite := hasFlag(flagOverwrite)
		err := Generate(apiFile, debug, overwrite)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func parseFilters() []string {

	filters := strings.Split(os.Getenv("FILTER"), ",")

	if len(os.Args) > 1 {
		args := os.Args[1:]
		for _, arg := range args {
			if strings.Contains(arg, fmt.Sprintf("%s=", paramFilter)) {
				filters2 := strings.Split(strings.Split(arg, "=")[1], ",")
				filters = append(filters, filters2...)
			}
		}
	}

	filtersClean := []string{}
	for _, filter := range filters {
		if len(filter) > 0 {
			filtersClean = append(filtersClean, filter)
		}
	}

	return filtersClean
}

func hasFlag(flagValue string) bool {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		for _, arg := range args {
			if strings.Trim(arg, "- ") == flagValue {
				return true
			}
		}
	}
	return false
}

func getBluezVersion() string {

	bluezVersion, err := util.GetGitVersion(docsDir)
	if err != nil {
		log.Fatal(err)
	}

	envBluezVersion := os.Getenv("BLUEZ_VERSION")
	if envBluezVersion != "" {
		bluezVersion = envBluezVersion
	}

	log.Infof("API %s", bluezVersion)
	return bluezVersion
}

func parseLogLevel() logrus.Level {
	logLevel := log.DebugLevel.String()
	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)
	return lvl
}

func Parse(filters []string, debug bool) error {

	api, err := gen.Parse(docsDir, filters, debug)
	if err != nil {
		log.Fatalf("Parse failed: %s", err)
		return err
	}

	apiFile := fmt.Sprintf("./bluez-%s.json", api.Version)
	log.Infof("Saving to %s\n", apiFile)
	err = api.Serialize(apiFile)
	if err != nil {
		log.Fatalf("Failed to serialize JSON: %s", err)
		return err
	}

	return nil
}

func Generate(filename string, debug bool, overwrite bool) error {

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

	err = generator.Generate(api, "./bluez", debug, overwrite)
	if err != nil {
		log.Fatalf("Generation failed: %s", err)
		return err
	}

	return nil
}

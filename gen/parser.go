package gen

import (
	"strings"

	"github.com/muka/go-bluetooth/gen/filters"
	"github.com/muka/go-bluetooth/gen/parser"
	"github.com/muka/go-bluetooth/gen/types"
	"github.com/muka/go-bluetooth/gen/util"
	log "github.com/sirupsen/logrus"
)

// Parse bluez DBus API docs
func Parse(docsDir string, filtersList []filters.Filter, debug bool) (BluezAPI, error) {
	files, err := util.ListFiles(docsDir)
	if err != nil {
		return BluezAPI{}, err
	}
	apis := make([]*types.ApiGroup, 0)
	for _, file := range files {

		keep := true
		if len(filtersList) > 0 {
			keep = false
			for _, filter1 := range filtersList {
				if filter1.Context != filters.FilterFile {
					continue
				}
				if strings.Contains(file, filter1.Value) {
					keep = true
					if debug {
						log.Debugf("[filter %s] Keep %s", filter1.Value, file)
					}
					break
				}
			}
		}

		if !keep {
			continue
		}

		apiGroupParser := parser.NewApiGroupParser(debug, filtersList)
		apiGroup, err := apiGroupParser.Parse(file)
		if err != nil {
			log.Errorf("Failed to load %s, skipped", file)
			continue
		}
		apis = append(apis, apiGroup)
	}

	version, err := util.GetGitVersion(docsDir)
	if err != nil {
		log.Errorf("Failed to parse version: %s", err)
	}

	return BluezAPI{
		Version: version,
		Api:     apis,
	}, nil
}

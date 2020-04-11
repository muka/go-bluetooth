package gen

import (
	"strings"

	"github.com/muka/go-bluetooth/gen/parser"
	"github.com/muka/go-bluetooth/gen/types"
	"github.com/muka/go-bluetooth/gen/util"
	log "github.com/sirupsen/logrus"
)

// Parse bluez DBus API docs
func Parse(docsDir string, filters []string, debug bool) (BluezAPI, error) {
	files, err := util.ListFiles(docsDir)
	if err != nil {
		return BluezAPI{}, err
	}
	apis := make([]*types.ApiGroup, 0)
	for _, file := range files {

		keep := true
		if len(filters) > 0 {
			keep = false
			for _, filter := range filters {
				if strings.Contains(file, filter) {
					keep = true
					if debug {
						log.Debugf("[filter %s] Keep %s", filter, file)
					}
					break
				}
			}
		}

		if !keep {
			continue
		}

		apiGroupParser := parser.NewApiGroupParser(debug)
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

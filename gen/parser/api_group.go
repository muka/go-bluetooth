package parser

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/muka/go-bluetooth/gen/filters"
	"github.com/muka/go-bluetooth/gen/types"
	"github.com/muka/go-bluetooth/gen/util"

	log "github.com/sirupsen/logrus"
)

type ApiGroupParser struct {
	model  *types.ApiGroup
	debug  bool
	filter []filters.Filter
}

// NewApiGroupParser parser for ApiGroup
func NewApiGroupParser(debug bool, filtersList []filters.Filter) ApiGroupParser {
	apiGroupParser := ApiGroupParser{
		debug:  debug,
		filter: filtersList,
		model: &types.ApiGroup{
			Api: make([]*types.Api, 0),
		},
	}
	return apiGroupParser
}

// Parse load a documentation file and parse the content
func (g *ApiGroupParser) Parse(srcFile string) (*types.ApiGroup, error) {

	var err error
	apiGroup := g.model

	if g.debug {
		log.Debugf("------------------- Parsing %s -------------------", srcFile)
	}

	apiGroup.FileName = filepath.Base(srcFile)
	apiGroup.Api = make([]*types.Api, 0)

	raw, err := util.ReadFile(srcFile)
	if err != nil {
		return apiGroup, err
	}

	// Split by sections eg
	// group name
	// ******
	// group description ...
	// api title
	// ======
	// api contents..

	// Apis
	re1 := regexp.MustCompile(`([A-Za-z0-1._ -]*)\n[=]+\n`)
	matches1 := re1.FindAllSubmatchIndex(raw, -1)

	if len(matches1) == 0 {
		return apiGroup, fmt.Errorf("%s: no service defined?", srcFile)
	}

	// split up
	groupText := raw[:matches1[0][0]]
	g.parseGroup(groupText)

	// log.Debugf("%d", matches1)

	slices := make([][]byte, 0)
	if len(matches1) == 1 {
		serviceRaw := raw[matches1[0][0]:]
		if len(serviceRaw) > 0 {
			slices = append(slices, serviceRaw)
		}
	} else {

		prevPos := -1
		for i := 0; i < len(matches1); i++ {

			if prevPos == -1 {
				prevPos = matches1[i][0]
				continue
			}

			currPos := matches1[i][0]
			serviceRaw := raw[prevPos:currPos]
			prevPos = currPos

			// log.Debugf("%s", serviceRaw)

			if len(serviceRaw) > 0 {
				slices = append(slices, serviceRaw)
			}

			// keep the last one
			if i == len(matches1)-1 {
				serviceRaw = raw[currPos:]
				slices = append(slices, serviceRaw)
			}

		}
	}

	for _, slice := range slices {
		apiParser := NewApiParser(g.debug, g.filter)
		api, err := apiParser.Parse(slice)
		if err != nil {
			return apiGroup, err
		}
		apiGroup.Api = append(apiGroup.Api, api)
	}

	return apiGroup, nil
}

// func (g *ApiParser) parseApi(raw []byte) (*types.Api, error) {
// 	apiParser := NewApiParser(g.debug, g.filter)
// 	return apiParser.Parse(raw)
// }

func (g *ApiGroupParser) parseGroup(raw []byte) {

	// Group Name
	re := regexp.MustCompile(`(.+)\n[*]+\n(.*)`)
	matches := re.FindAllSubmatchIndex(raw, -1)

	// log.Debugf("\nRAW\n%s\n\n/RAW\n", raw)
	// for _, m1 := range matches {
	// 	// for _, m := range m1 {
	// 	log.Debugf("> %v", m1)
	// 	// }
	// }

	g.model.Name = string(raw[matches[0][2]:matches[0][3]])
	g.model.Description = string(raw[matches[0][1]+1:])

	if g.debug {
		log.Debugf("* %s", g.model.Name)
	}
}

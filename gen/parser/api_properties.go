package parser

import (
	"regexp"

	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

func (g *ApiParser) ParseProperties(raw []byte) ([]*types.Property, error) {

	var err error = nil
	props := make([]*types.Property, 0)
	slices := make([][]byte, 0)

	re := regexp.MustCompile(`(?s)\nProperties(.+)\n\n?(Filters|)?[ \t]?`)
	matches1 := re.FindSubmatch(raw)

	if len(matches1) == 0 {
		return props, err
	}

	for _, propsRaw := range matches1[1:] {

		// string Modalias [readonly, optional]
		re1 := regexp.MustCompile(`(?s)[ \t]*` + propBaseRegexp + `.*?\n`)
		matches2 := re1.FindAllSubmatchIndex(propsRaw, -1)

		// log.Debugf("1*** %d", matches2)

		// if len(matches2) == 0 {
		// re1 := regexp.MustCompile(`[ \t]*(bool|byte|string|uint|dict|array.*) ([A-Za-z0-9_]+?)( ?) *\n`)
		// matches2 := re1.FindAllSubmatchIndex(propsRaw, -1)
		// }

		// log.Debugf("2 *** %d", matches2)

		if len(matches2) == 1 {
			if len(propsRaw) > 0 {
				// log.Debugf("ADD single *** %s", propsRaw)
				slices = append(slices, propsRaw)
			}
		} else {
			prevPos := -1
			for i := 0; i < len(matches2); i++ {

				if prevPos == -1 {
					prevPos = matches2[i][0]
					continue
				}

				nextPos := matches2[i][0]
				propRaw := propsRaw[prevPos:nextPos]
				prevPos = nextPos

				if len(propRaw) > 0 {
					slices = append(slices, propRaw)
				}

				// keep the last one
				lastItem := len(matches2) - 1
				if i == lastItem {
					propRaw = propsRaw[matches2[lastItem][0]:]
					if len(propRaw) > 0 {
						slices = append(slices, propRaw)
					}
				}
			}
		}
	}

	if g.debug {
		log.Debug("\tProperties:")
	}

	for _, propRaw := range slices {
		propertyParser := NewPropertyParser(g.debug)
		prop, err := propertyParser.Parse(propRaw)
		if err != nil {
			log.Warnf("Skipped property: %s", err)
			continue
		}
		props = append(props, prop)
	}

	return props, nil
}

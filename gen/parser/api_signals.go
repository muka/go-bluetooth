package parser

import (
	"regexp"

	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

func (g *ApiParser) ParseSignals(raw []byte) ([]*types.Method, error) {

	var err error
	methods := make([]*types.Method, 0)
	slices := make([][]byte, 0)

	re := regexp.MustCompile(`(?s)Signals(.+)\n\nProperties`)
	matches1 := re.FindSubmatch(raw)

	if len(matches1) == 0 {
		return methods, err
	}

	// if len(matches1) == 0 {
	// 	re = regexp.MustCompile(`(?s)[ \t\n]+(.+)`)
	// 	matches1 = re.FindSubmatch(raw)
	// 	if len(matches1) == 1 {
	// 		matches1 = append(matches1, matches1[0])
	// 	}
	// }

	// log.Debugf("matches1 %s", matches1[1:])
	// log.Debugf("%s", matches1)

	for _, methodsRaw := range matches1[1:] {

		re1 := regexp.MustCompile(`[ \t]*?(.*?)? ?([^ ]+)\(([^)]+?)?\) ?(.*)`)
		matches2 := re1.FindAllSubmatchIndex(methodsRaw, -1)

		if len(matches2) == 1 {
			if len(methodsRaw) > 0 {
				slices = append(slices, methodsRaw)
			}
		} else {
			prevPos := -1
			for i := 0; i < len(matches2); i++ {

				if prevPos == -1 {
					prevPos = matches2[i][0]
					continue
				}

				nextPos := matches2[i][0]
				methodRaw := methodsRaw[prevPos:nextPos]
				prevPos = nextPos

				if len(methodRaw) > 0 {
					slices = append(slices, methodRaw)
				}

				// keep the last one
				lastItem := len(matches2) - 1
				if i == lastItem {
					methodRaw = methodsRaw[matches2[lastItem][0]:]
					if len(methodRaw) > 0 {
						slices = append(slices, methodRaw)
					}
				}
			}
		}
	}

	if g.debug {
		log.Debug("\nSignals:")
	}
	for _, methodRaw := range slices {
		methodParser := NewMethodParser(g.debug)
		method, err := methodParser.Parse(methodRaw)
		if err != nil {
			log.Debugf("Skip signal: %s", err)
			continue
		}
		methods = append(methods, method)
	}

	return methods, nil
}

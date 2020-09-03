package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/gen/filters"
	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

func (g *ApiParser) ParseMethods(raw []byte) ([]*types.Method, error) {

	var err error = nil
	methods := make([]*types.Method, 0)
	slices := make([][]byte, 0)

	hasMethods := true

	re := regexp.MustCompile(`(?sm)Methods:?\n?(\t+.*?)\n(?:Properties|Filter|Signals):?\n?`)
	matches1 := re.FindAllSubmatch(raw, -1)

	if len(matches1) == 0 {
		onlyPropsRegex := regexp.MustCompile(`(?smi)(Properties|Filter|Signals):?\n?\t+`)
		onlyProps := onlyPropsRegex.Match(raw)
		if !onlyProps {
			matches1 = append(matches1, [][]byte{nil, raw})
			log.Tracef("[%s] No methods wrapper matched, take all text", g.model.Title)
		} else {
			hasMethods = false
		}
		// log.Debugf("matches1 %s", matches1[1:])
	}

	for _, matches1list := range matches1 {

		methodsRaw := matches1list[1]

		// crop from [Methods:]
		mRe := regexp.MustCompile("Methods:?\n?")
		methodsPos := mRe.FindIndex(methodsRaw)

		if len(methodsPos) > 0 {
			methodsRaw = methodsRaw[methodsPos[1]:]
		}

		// remove initial \n chars that are not tabs
		rmLeadingRegex := regexp.MustCompile(`(?smi)^[^\t]*`)
		methodsRaw = rmLeadingRegex.ReplaceAll(methodsRaw, []byte{})
		// methodsRaw = bytes.TrimLeft(methodsRaw, "\n")

		// log.Debugf("methodsRaw\n%s", methodsRaw)

		tabLength := 0
		for {
			if methodsRaw[tabLength] == '\t' {
				tabLength++
				continue
			}
			break
		}

		// at least one tab
		if tabLength == 0 {
			tabLength++
		}

		log.Tracef("[%s] Methods tab spacing is %d", g.model.Title, tabLength)

		// re1 := regexp.MustCompile(`[ \t]*?(.*?)? ?([^ ]+)\(([^)]+?)?\) ?(.*)`)
		re1 := regexp.MustCompile(fmt.Sprintf(`(?ms)(^[\t]{%d}[^\t]+)`, tabLength))
		matches2 := re1.FindAllIndex(methodsRaw, -1)

		// log.Debug(strings.ReplaceAll(string(methodsRaw), "\t", "->"))
		// os.Exit(1)

		// log.Debugf("%v", matches2)
		// os.Exit(1)

		// take by method line to the next method line, including anything in the middle
		// if there is just one method, take all
		if len(matches2) == 1 {
			if len(methodsRaw) > 0 {
				slices = append(slices, methodsRaw)
			}
		} else {

			if len(matches2) == 0 {
				g.log("No methods found")
			}

			lastValue := []byte{}
			prevPos := -1
			for i := 0; i < len(matches2); i++ {

				if prevPos == -1 {
					prevPos = matches2[i][0]
					continue
				}

				nextPos := matches2[i][0]
				methodRaw := methodsRaw[prevPos:nextPos]

				if len(lastValue) > 0 {
					methodRaw = append(lastValue, methodRaw...)
					lastValue = []byte{}
				}

				prevPos = nextPos

				// obex api docs split return args from function name keeping same tabs,
				// aggregate if methodRaw seems short
				if len(methodRaw) < 40 {
					lastValue = methodRaw
					continue
				}

				// log.Tracef("raw method: %s", methodRaw)

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

	if g.debug && hasMethods {
		log.Debug("\tMethods:")
	}

	for _, methodRaw := range slices {
		methodParser := NewMethodParser(g.debug)
		method, err := methodParser.Parse(methodRaw)

		if err != nil {
			log.Traceln("------")
			log.Warnf("[%s] method parse error: %s", g.model.Title, err)
			log.Tracef("[%s]: %s", g.model.Title, methodRaw)
			log.Traceln("------")
			continue
		}

		// apply filters
		if len(g.filter) > 0 {

			skipMethod := false

			for _, f := range g.filter {

				if f.Context != filters.FilterMethod {
					continue
				}
				skipMethod = !strings.Contains(strings.ToLower(method.Name), strings.ToLower(f.Value))
			}

			if !skipMethod {
				methods = append(methods, method)
				// log.Debugf("Keep filtered method %s", method.Name)
			} else {
				// log.Debugf("Skip filtered method %s", method.Name)
			}

		}

		// keep all
		methods = append(methods, method)
	}

	if hasMethods && len(methods) == 0 {
		log.Warnf("%s: No methods found", g.model.Title)
	}

	return methods, err
}

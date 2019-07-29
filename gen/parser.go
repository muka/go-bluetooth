package gen

import (
	"fmt"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func (g *ApiGroup) Parse(srcFile string) error {

	log.Debugf("------------------- Parsing %s -------------------", srcFile)

	raw, err := ReadFile(srcFile)
	if err != nil {
		return err
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
		return fmt.Errorf("%s: no service defined?", srcFile)
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
		api := g.parseApi(slice)
		g.Api = append(g.Api, api)
	}

	return nil
}

func NewApiGroup(srcFile string) (ApiGroup, error) {
	g := ApiGroup{
		FileName: filepath.Base(srcFile),
		Api:      make([]Api, 0),
		debug:    true,
	}
	err := g.Parse(srcFile)
	return g, err
}

// Parse bluez DBus API docs and generate go code stub
func Parse(src string) []ApiGroup {
	files := ListFiles(src + "/doc")
	apis := []ApiGroup{}
	for _, file := range files {
		apiGroup, err := NewApiGroup(file)
		if err != nil {
			log.Errorf("Failed to load %s", file)
			continue
		}
		apis = append(apis, apiGroup)
	}
	return apis
}

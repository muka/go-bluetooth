package gen

import (
	"fmt"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type Flag int

const (
	FlagReadOnly     Flag = 1
	FlagWriteOnly    Flag = iota
	FlagReadWrite    Flag = iota
	FlagExperimental Flag = iota
)

type Arg struct {
	Type string
	Name string
}

type Method struct {
	Name       string
	ReturnType string
	Args       []Arg
	Errors     []string
	Docs       string
}

type Property struct {
	Name  string
	Type  string
	Docs  string
	Flags []Flag
}

type ApiGroup struct {
	FileName    string
	Name        string
	Description string
	Api         []Api
}

type Api struct {
	Title       string
	Description string
	Service     string
	Interface   string
	ObjectPath  string
	Methods     []Method
	Properties  []Property
}

func (g *ApiGroup) Parse(srcFile string) error {

	log.Debugf("------------------- Parsing %s -------------------", srcFile)

	raw, err := readFile(srcFile)
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
		g.parseApi(slice)
	}

	return nil
}

func NewApiGroup(srcFile string) (ApiGroup, error) {
	g := ApiGroup{
		FileName: filepath.Base(srcFile),
		Api:      make([]Api, 0),
	}
	err := g.Parse(srcFile)
	return g, err
}

// Parse bluez DBus API docs and generate go code stub
func Parse(src string) []ApiGroup {
	apis := listFiles(src + "/doc")
	return apis
}

func Generate(outDir string) {

}

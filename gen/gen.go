package gen

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Method struct {
	Name       string
	ReturnType string
	Args       []string
	Errors     []string
	Docs       string
}

type ApiGroup struct {
	FileName    string
	Name        string
	Description string
	Api         []Api
}

type Api struct {
	Name       string
	Service    string
	Interface  string
	ObjectPath string
	Methods    []Method
}

func (g *ApiGroup) read(srcFile string) ([]byte, error) {
	file, err := os.Open(srcFile)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func (g *ApiGroup) parseApi(raw []byte) {

}

func (g *ApiGroup) parseGroup(raw []byte) {

	// Group Name
	re := regexp.MustCompile(`(.+)\n[*]+\n(.*)`)
	matches := re.FindAllSubmatchIndex(raw, -1)

	// log.Debugf("\nRAW\n%s\n\n/RAW\n", raw)
	// for _, m1 := range matches {
	// 	// for _, m := range m1 {
	// 	log.Debugf("> %v", m1)
	// 	// }
	// }

	g.Name = string(raw[matches[0][2]:matches[0][3]])
	g.Description = string(raw[matches[0][1]+1:])

	log.Debugf("* %s", g.Name)
	log.Debugf("*\t%s", g.Description)

}

func (g *ApiGroup) Parse(srcFile string) error {

	log.Debugf("Parsing %s", srcFile)

	raw, err := g.read(srcFile)
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
	re1 := regexp.MustCompile(`(?P<service>[A-Za-z0-1._ -]*)\n[=]+\n`)
	matches1 := re1.FindAllSubmatchIndex(raw, -1)

	if len(matches1) == 0 {
		log.Warnf("doc has no service defined?")
		return nil
	}

	// split up
	groupText := raw[:matches1[0][0]]
	g.parseGroup(groupText)

	slices := make([][]byte, len(matches1))
	prevPos := 0
	for i := 0; i < len(matches1); i++ {
		if prevPos > 0 {
			slices[i] = raw[prevPos:matches1[i][0]]
		}
		prevPos = matches1[i][0]
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

func listFiles(dir string) []ApiGroup {

	list := make([]ApiGroup, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "-api.txt") {
			apiGroup, err := NewApiGroup(path)
			if err != nil {
				log.Errorf("Failed to load %s", path)
				return nil
			}
			list = append(list, apiGroup)
		}
		return nil
	})

	if err != nil {
		log.Errorf("Failed to list files: %s", err)
	}

	return list
}

func Generate(outDir string) {

}

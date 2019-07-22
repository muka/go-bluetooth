package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
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
	Name     string
	Type     string
	Docs     string
	readable bool
	writable bool
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

func (g *ApiGroup) parseProperty(raw []byte) Property {

	// log.Debugf("prop raw -> %s", raw)

	re1 := regexp.MustCompile(`[ \t]*([A-Za-z}{]+) (.+?) \[(.*)\]\n((?s).+)`)
	matches2 := re1.FindAllSubmatch(raw, -1)
	// log.Debugf("m1 %s", matches2)
	if len(matches2) == 0 {
		re1 = regexp.MustCompile(`[ \t]*([A-Za-z}{]+) (.+)( ?)\n((?s).+)`)
		matches2 = re1.FindAllSubmatch(raw, -1)
	}

	// log.Debugf("m2 %s", matches2)

	docs := string(matches2[0][4])
	docs = strings.Replace(docs, " \t\n", "", -1)
	docs = strings.Trim(docs, " \t\n")

	name := string(matches2[0][2])
	name = strings.Replace(name, " \t\n", "", -1)

	p := Property{
		Type: string(matches2[0][1]),
		Name: name,
		// : string(matches2[0][3]),
		Docs: docs,
	}
	// log.Debugf("\t %s %s", p.Type, p.Name)
	log.Debugf("\t - %s %s", p.Type, p.Name)
	return p
}

func (g *ApiGroup) parseProperties(raw []byte) []Property {

	props := make([]Property, 0)
	slices := make([][]byte, 0)

	re := regexp.MustCompile(`(?s)\nProperties(.+)\n\n`)
	matches1 := re.FindSubmatch(raw)

	for _, propsRaw := range matches1[1:] {

		// string Modalias [readonly, optional]
		re1 := regexp.MustCompile(`[ \t]*(byte|string|uint|dict|array.*)[ \t](.+) \[(.*)\] *\n`)
		matches2 := re1.FindAllSubmatchIndex(propsRaw, -1)

		// log.Debugf("1*** %d", matches2)

		if len(matches2) == 0 {
			re1 = regexp.MustCompile(`[ \t]*(byte|string|uint|dict|array.*)[ \t](.+)( ?) *\n`)
			matches2 = re1.FindAllSubmatchIndex(propsRaw, -1)
		}

		// log.Debugf("2*** %d", matches2)

		if len(matches2) == 1 {
			if len(propsRaw) > 0 {
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

	log.Debug("\tProperties:")
	for _, propRaw := range slices {
		prop := g.parseProperty(propRaw)
		props = append(props, prop)
	}

	return props
}

func (g *ApiGroup) parseMethods(raw []byte) []Method {

	methods := make([]Method, 0)
	slices := make([][]byte, 0)

	re := regexp.MustCompile(`(?s)Methods(.+)\n\nProperties`)
	matches1 := re.FindSubmatch(raw)

	// handle agent-api.txt case
	if len(matches1) == 0 {
		re = regexp.MustCompile(`(?s)[ \t\n]+(.+)`)
		matches1 = re.FindSubmatch(raw)
		if len(matches1) == 1 {
			matches1 = append(matches1, matches1[0])
		}
	}

	// log.Debugf("matches1 %s", matches1[1:])
	// log.Debugf("%s", matches1)

	for _, methodsRaw := range matches1[1:] {

		re1 := regexp.MustCompile(`[ \t]*?(.*?) ?(.+)\(([^)]+?)?\) ?(.*)`)
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

	log.Debug("\tMethods:")
	for _, methodRaw := range slices {
		method := g.parseMethod(methodRaw)
		methods = append(methods, method)
	}

	return methods
}

func (g *ApiGroup) parseMethod(raw []byte) Method {

	method := Method{}
	// log.Debugf("%s", raw)

	re := regexp.MustCompile(`[ \t]*(.*?) ?(\w+)\(([^)]*)\) ?(.*?)\n((?s).+)`)
	matches1 := re.FindAllSubmatch(raw, -1)

	// log.Debugf("matches1 %s", matches1)
	for _, matches2 := range matches1 {

		rtype := string(matches2[1])
		if len(rtype) > 7 && rtype[:7] == "Methods" {
			rtype = rtype[7:]
		}
		method.ReturnType = strings.Trim(rtype, " \t")

		name := string(matches2[2])
		method.Name = strings.Trim(name, " \t")

		args := []Arg{}
		if len(matches2[3]) > 0 {
			argslist := strings.Split(string(matches2[3]), ",")
			for _, arg := range argslist {
				arg = strings.Trim(arg, " ")
				argsparts := strings.Split(arg, " ")
				arg := Arg{
					Type: strings.Trim(argsparts[0], " \t\n"),
					Name: argsparts[1],
				}
				args = append(args, arg)
			}
		}
		method.Args = args
		method.Docs = string(matches2[5])
	}

	log.Debugf("\t - %s %s(%s)", method.ReturnType, method.Name, method.Args)

	return method
}

func (g *ApiGroup) parseApi(raw []byte) {

	api := Api{}

	// title & description
	re := regexp.MustCompile(`(?s)(.+)\n[=]+\n(.*)\nService|Interface *`)
	matches := re.FindSubmatchIndex(raw)

	api.Title = string(raw[matches[2]:matches[3]])
	api.Description = string(raw[matches[4]:matches[5]])

	log.Debugf("= %s", api.Title)

	raw = raw[matches[5]:]

	// service interface object
	re = regexp.MustCompile(`Service[ \t]*(.+)\nInterface[ \t]*(.+)\nObject path[ \t]*(.+)\n`)
	matches = re.FindSubmatchIndex(raw)

	api.Service = string(raw[matches[2]:matches[3]])
	api.Interface = string(raw[matches[4]:matches[5]])
	api.ObjectPath = string(raw[matches[6]:matches[7]])

	if strings.Contains(api.ObjectPath, "freely") {
		api.ObjectPath = ""
	}

	log.Debugf("\tService %s", api.Service)
	log.Debugf("\tInterface %s", api.Interface)
	log.Debugf("\tObjectPath %s", api.ObjectPath)

	raw = raw[matches[7]:]

	api.Methods = g.parseMethods(raw)
	api.Properties = g.parseProperties(raw)

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
	// log.Debugf("*\t%s", g.Description)

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
		return fmt.Errorf("%s: no service defined?", srcFile)
	}

	// split up
	groupText := raw[:matches1[0][0]]
	g.parseGroup(groupText)

	slices := make([][]byte, 0)
	prevPos := -1
	for i := 0; i < len(matches1); i++ {

		if prevPos == -1 {
			prevPos = matches1[i][0]
			continue
		}

		currPos := matches1[i][0]
		serviceRaw := raw[prevPos:currPos]
		prevPos = currPos

		if len(serviceRaw) > 0 {
			slices = append(slices, serviceRaw)
		}

		// keep the last one
		if i == len(matches1)-1 {
			serviceRaw = raw[currPos:]
			slices = append(slices, serviceRaw)
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

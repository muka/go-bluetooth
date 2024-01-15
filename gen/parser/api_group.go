package parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

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
	if strings.HasSuffix(srcFile, ".txt") { // nolint: gocritic
		return g.parseTXT(srcFile)
	} else if strings.HasSuffix(srcFile, ".rst") {
		return g.parseRST(srcFile)
	} else {
		log.Errorf("Unknown file type for %s", srcFile)
		return nil, fmt.Errorf("Unknown file type.")
	}
}

func (g *ApiGroupParser) getSection(raw, sectionName string, divider rune) (out string, err error) {
	pattern := fmt.Sprintf("(?ms)%s\\n%c+\\n(.*?)\\n((Methods|Signals|Properties)\\n[-`]+\\n|\\z)", sectionName, divider)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(raw)
	if len(matches) > 0 {
		return strings.Trim(matches[1], " \t\n"), nil
	}

	return "", nil
}

func (g *ApiGroupParser) parseFlags(raw string) (flags []types.Flag, err error) {
	for _, f := range strings.Split(raw, ", ") {
		var flag types.Flag = 0
		switch f {
		case "readonly":
			fallthrough
		case "read-only":
			// log.Printf("f: %v", f)
			flag = types.FlagReadOnly
		case "writeonly":
			fallthrough
		case "write-only":
			// log.Printf("f: %v", f)
			flag = types.FlagWriteOnly
		case "readwrite":
			fallthrough
		case "read-write":
			fallthrough
		case "read/write":
			// log.Printf("f: %v", f)
			flag = types.FlagReadWrite
		case "experimental":
			fallthrough
		case "Experimental":
			// log.Printf("f: %v", f)
			flag = types.FlagExperimental
		case "optional":
			// log.Printf("f: %v", f)
			flag = types.FlagOptional
		default:
			log.Warnf("Unknown flag %s", f)
		}
		if flag != 0 {
			flags = append(flags, flag)
		}
	}

	return flags, nil
}

func (g *ApiGroupParser) parseRST(srcFile string) (*types.ApiGroup, error) {
	var err error
	apiGroup := g.model

	if g.debug {
		log.Debugf("------------------- Parsing %s -------------------", srcFile)
	}

	apiGroup.FileName = filepath.Base(srcFile)
	apiGroup.Api = make([]*types.Api, 0)

	rawBytes, err := util.ReadFile(srcFile)
	if err != nil {
		return apiGroup, err
	}

	raw := string(rawBytes)

	re := regexp.MustCompile(`-+\n([^\n]+)\n-+\n`)
	matches := re.FindStringSubmatch(string(raw))

	g.model.Name = matches[1]

	section, err := g.getSection(raw, "Description", '=')
	if err != nil {
		return apiGroup, err
	} else if section != "" {
		// log.Printf(" RST Description: %+#v", section)
		g.model.Description = strings.Trim(section, " \t\n")
	}

	api := &types.Api{}

	section, err = g.getSection(raw, "Interface", '=')
	if err != nil {
		return apiGroup, err
	} else if section != "" {
		re = regexp.MustCompile(`[:;]?Service[:;]?\s+(.*?)\n`)
		matches = re.FindStringSubmatch(section)
		api.Service = matches[1]
		// log.Printf("Service Match: %+#v", matches)

		re = regexp.MustCompile(`[:;]?Interface[:;]?\s+([^ \n]+)`)
		matches = re.FindStringSubmatch(section)
		api.Interface = matches[1]
		// log.Printf("Interface Match: %+#v", matches)

		re = regexp.MustCompile(`[:;]?Object path[:;]?\s+([^\n]+)`)
		matches = re.FindStringSubmatch(section)
		api.ObjectPath = matches[1]
		// log.Printf("ObjectPath Match: %+#v", matches)

		api.Title = g.model.Name
	}

	section, err = g.getSection(raw, "Methods", '-')
	switch {
	case err != nil:
		return apiGroup, err
	case section != "":
		// log.Printf("Methods Section: %s", section)

		re = regexp.MustCompile("(?ms)^(([^\\s`][^\n]*)\\s(\\w*)\\(([^\n)]*)\\))([^\\n\\r]*)\\n`+\\n(\\n(\\t[^\\n]*[\\n\\r]+|[\\n\\r]+)+)")
		methodMatches := re.FindAllStringSubmatch(section, -1)
		for _, x := range methodMatches {
			method := &types.Method{}
			method.Name = x[3]
			method.ReturnType = x[2]
			method.Docs = strings.Trim(x[6], " \t\n")
			// log.Printf(" RST Method: %s - %s - %s, %+#v", x[1], x[2], x[3], x)

			if x[4] != "" {
				for _, arg := range strings.Split(x[4], ", ") {
					// log.Printf("  RST Method Arg: %+#v", arg)
					v := strings.Split(arg, " ")
					if len(v) == 1 { // FIXME: This is a horrible hack to deal with org.bluez.Profile.rst -> NewConnection missing a type for the fd argument.
						method.Args = append(method.Args, types.Arg{Name: arg, Type: "int32"})
					} else {
						method.Args = append(method.Args, types.Arg{Name: v[1], Type: v[0]})
					}
				}
			} else {
				method.Args = []types.Arg{}
			}

			re = regexp.MustCompile("Possible errors:\n\n((?:\t:.*:\n)+)")
			errors := re.FindStringSubmatch(method.Docs)
			// log.Printf(" RST Errors: %+#v", errors)
			if errors != nil {
				re = regexp.MustCompile("\t:(.*):\n")
				errs := re.FindAllStringSubmatch(errors[1], -1)
				// log.Printf("  RST Errors: %#v", errs)
				for _, x := range errs {
					method.Errors = append(method.Errors, x[1])
				}
			} else {
				method.Errors = []string{}
			}

			// log.Printf("  RST Method: %+#v", method)
			api.Methods = append(api.Methods, method)
		}
	default:
		api.Methods = []*types.Method{}
	}

	api.Signals = []*types.Method{}

	section, err = g.getSection(raw, "Properties", '-')
	switch {
	case err != nil:
		return apiGroup, err
	case section != "":
		// log.Printf("Properties Section: %s", section)

		// re := regexp.MustCompile("(?:\\n|^)((.*?) (\\w+)(?: \\[([a-z, -]+)\\])?)\\n`+\\n(\\n(\\t[^\\n]*[\\n\\r]+|[\\n\\r]+)+)")
		// re := regexp.MustCompile("(?:\\n|^)((.*?) (\\w+)(?: \\[([a-z, -]+)\\])?)\\n`+\\n((\\n(\\t[^\\n]*[\\n\\r]+|[\\n\\r]+))+)")
		re := regexp.MustCompile(
			"(?ms:^)((.*?)" +
				" (\\w+)" +
				"(?: \\[([a-z, -]+)\\])?" +
				")(?: \\(Default:.*\\))?\\n`+\\n+" +

				"(" +
				"(?:" +
				"\\t.+" +
				"|" +
				"\\n" +
				")*" +
				")")
		// re := regexp.MustCompile("(?:\\n|^)((.*?) (\\w+)(?: \\[([a-z, -]+)\\])?)\\n`+")
		matches := re.FindAllStringSubmatch(section, -1)
		// log.Printf(" RST Properties: %+#v", matches)

		for _, x := range matches {
			property := &types.Property{}
			property.Name = x[3]
			property.Type = x[2]
			property.Docs = strings.Trim(x[5], " \t\n")
			// property.Docs = x[5]
			// log.Printf("  RST Property: %s - %s - %s - %s, %+#v", x[1], x[2], x[3], x[4], x)

			if x[4] != "" {
				property.Flags, err = g.parseFlags(x[4])
				if err != nil {
					return apiGroup, err
				}
			}

			// log.Printf("  RST Property: %+#v", property)
			api.Properties = append(api.Properties, property)
		}
	default:
		api.Properties = []*types.Property{}
	}

	// One API per file for the RST files.
	apiGroup.Api = append(apiGroup.Api, api)
	return apiGroup, nil
}

func (g *ApiGroupParser) parseTXT(srcFile string) (*types.ApiGroup, error) {
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

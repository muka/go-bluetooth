package parser

import (
	"errors"
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

const propBaseRegexp = `(bool|boolean|byte|string|[i|I]nt16|[U|u]int16|uint16_t|uint32|dict|object|array\{.*?) ([A-Z].+?)`

type PropertyParser struct {
	model *types.Property
	debug bool
}

// NewPropertyParser
func NewPropertyParser(debug bool) PropertyParser {
	p := PropertyParser{
		model: new(types.Property),
		debug: debug,
	}
	return p
}

func (g *PropertyParser) Parse(raw []byte) (*types.Property, error) {

	var err error
	property := g.model
	// log.Debugf("prop raw -> %s", raw)

	re1 := regexp.MustCompile(`[ \t]*?` + propBaseRegexp + `( \[[^\]]*\].*)?\n((?s).+)`)
	matches2 := re1.FindAllSubmatch(raw, -1)

	// log.Warnf("m1 %s", matches2)

	if len(matches2) == 0 || len(matches2[0]) == 1 {
		re1 = regexp.MustCompile(`[ \t]*?` + propBaseRegexp + `\n((?s).+)`)
		matches2 = re1.FindAllSubmatch(raw, -1)
		// log.Warnf("m2 %s", matches2)
	}

	if len(matches2) == 0 {
		log.Debugf("prop raw -> %s", raw)
		return property, errors.New("no property found")
	}

	flags := []types.Flag{}
	flagListRaw := string(matches2[0][3])
	flagList := strings.Split(strings.Trim(flagListRaw, "[] "), ",")

	for _, f := range flagList {

		// track server-only flags for gatt API
		if strings.Contains(f, "Server Only") {
			flags = append(flags, types.FlagServerOnly)
		}

		// int16 Handle [read-write, optional] (Server Only)
		if strings.Contains(f, "]") {
			f = strings.Split(f, "]")[0]
		}

		f = strings.Trim(f, " []")
		if f != "" {
			var flag types.Flag = 0
			switch f {
			case "readonly":
			case "read-only":
				flag = types.FlagReadOnly
			case "writeonly":
			case "write-only":
				flag = types.FlagWriteOnly
			case "readwrite":
			case "read-write":
			case "read/write":
				flag = types.FlagReadWrite
			case "experimental":
			case "Experimental":
				flag = types.FlagExperimental
			case "optional":
				flag = types.FlagOptional
			default:
				log.Warnf("Unknown flag %s", f)
			}

			if flag > 0 {
				flags = append(flags, flag)
			}
		}
	}

	docs := string(matches2[0][4])
	docs = strings.Replace(docs, " \t\n", "", -1)
	docs = strings.Trim(docs, " \t\n")

	name := string(matches2[0][2])

	if strings.Contains(name, "optional") {
		name = strings.Replace(name, " (optional)", "", -1)
		docs = "(optional) " + docs
		flags = append(flags, types.FlagOptional)
	}

	name = strings.Replace(name, " \t\n", "", -1)

	// theese bastards fucks up with properties names
	if nameParts := strings.Split(name, " "); len(nameParts) > 1 {
		name = nameParts[0]
	}

	property.Type = string(matches2[0][1])
	property.Name = name
	property.Flags = flags
	property.Docs = docs

	if g.debug {
		log.Debugf("\t - %s", property)
	}
	return property, err
}

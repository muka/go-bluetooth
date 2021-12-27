package parser

import (
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/gen/filters"
	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

type ApiParser struct {
	model  *types.Api
	debug  bool
	filter []filters.Filter
}

// NewApiParser parser for Api
func NewApiParser(debug bool, filter []filters.Filter) ApiParser {
	parser := ApiParser{
		filter: filter,
		debug:  debug,
		model:  new(types.Api),
	}
	return parser
}

func (g *ApiParser) log(msg string) {
	log.Debugf("%s: %s", g.model.Title, msg)
}

func (g *ApiParser) Parse(raw []byte) (*types.Api, error) {

	var err error
	api := g.model

	// title & description
	re := regexp.MustCompile(`(?s)(.+)\n[=]+\n?(.*)\nService|Interface *`)
	matches := re.FindSubmatchIndex(raw)

	api.Title = string(raw[matches[2]:matches[3]])
	api.Description = string(raw[matches[4]:matches[5]])

	if g.debug {
		log.Debugf("= %s", api.Title)
	}

	if len(g.filter) > 0 {
		skipItem := false
		for _, filter := range g.filter {
			if filter.Context != filters.FilterApi {
				continue
			}
			skipItem = !strings.Contains(
				strings.ToLower(api.Title), strings.ToLower(filter.Value))
		}
		if skipItem {
			log.Debugf("Skip filtered API %s", api.Title)
			return nil, nil
		} else {
			log.Debugf("Keep filtered API %s", api.Title)
		}

	}

	raw = raw[matches[5]:]

	// service interface object
	re = regexp.MustCompile(`Service[ \t]*((?s).+)\nInterface[ \t]*((?s).+)\nObject path[ \t]*((?s).+?)\n\n`)
	matches = re.FindSubmatchIndex(raw)

	g.model = api
	api.Service = string(raw[matches[2]:matches[3]])
	api.Interface = strings.Replace(string(raw[matches[4]:matches[5]]), " [experimental]", "", -1)
	api.ObjectPath = string(raw[matches[6]:matches[7]])

	if g.debug {
		log.Debugf("\tService %s", api.Service)
		log.Debugf("\tInterface %s", api.Interface)
		log.Debugf("\tObjectPath %s", api.ObjectPath)
	}

	raw = raw[matches[7]:]

	methods, err := g.ParseMethods(raw)
	if err != nil {
		return api, err
	}
	api.Methods = methods

	properties, err := g.ParseProperties(raw)
	if err != nil {
		return api, err
	}
	api.Properties = properties

	signals, err := g.ParseSignals(raw)
	if err != nil {
		return api, err
	}
	api.Signals = signals

	return api, nil
}

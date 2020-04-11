package parser

import (
	"regexp"

	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

type ApiParser struct {
	model *types.Api
	debug bool
}

// NewApiParser parser for Api
func NewApiParser(debug bool) ApiParser {
	parser := ApiParser{
		debug: debug,
		model: new(types.Api),
	}
	return parser
}

func (g *ApiParser) Parse(raw []byte) (*types.Api, error) {

	var err error = nil
	api := g.model

	// title & description
	re := regexp.MustCompile(`(?s)(.+)\n[=]+\n?(.*)\nService|Interface *`)
	matches := re.FindSubmatchIndex(raw)

	api.Title = string(raw[matches[2]:matches[3]])
	api.Description = string(raw[matches[4]:matches[5]])

	if g.debug {
		log.Debugf("= %s", api.Title)
	}

	raw = raw[matches[5]:]

	// service interface object
	re = regexp.MustCompile(`Service[ \t]*((?s).+)\nInterface[ \t]*((?s).+)\nObject path[ \t]*((?s).+?)\n\n`)
	matches = re.FindSubmatchIndex(raw)

	api.Service = string(raw[matches[2]:matches[3]])
	api.Interface = string(raw[matches[4]:matches[5]])
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

	g.model = api

	return api, nil
}

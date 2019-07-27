package gen

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

func (g *ApiGroup) parseApi(raw []byte) Api {

	api := Api{}

	// title & description
	re := regexp.MustCompile(`(?s)(.+)\n[=]+\n(.*)\nService|Interface *`)
	matches := re.FindSubmatchIndex(raw)

	api.Title = string(raw[matches[2]:matches[3]])
	api.Description = string(raw[matches[4]:matches[5]])

	log.Infof("= %s", api.Title)

	raw = raw[matches[5]:]

	// service interface object
	re = regexp.MustCompile(`Service[ \t]*((?s).+)\nInterface[ \t]*((?s).+)\nObject path[ \t]*((?s).+?)\n\n`)
	matches = re.FindSubmatchIndex(raw)

	// log.Debugf("%d", matches)

	api.Service = string(raw[matches[2]:matches[3]])
	api.Interface = string(raw[matches[4]:matches[5]])
	api.ObjectPath = string(raw[matches[6]:matches[7]])

	if g.debug {
		log.Debugf("\tService %s", api.Service)
		log.Debugf("\tInterface %s", api.Interface)
		log.Debugf("\tObjectPath %s", api.ObjectPath)
	}

	raw = raw[matches[7]:]
	api.Methods = g.parseMethods(raw)
	api.Signals = g.parseSignals(raw)
	api.Properties = g.parseProperties(raw)

	return api
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

	log.Infof("* %s", g.Name)

}

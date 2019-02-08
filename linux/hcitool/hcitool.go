package hcitool

import (
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/linux"
)

type HcitoolDev struct {
	ID      string
	Address string
}

// GetAdapter Return an adapter using hcitool as backend
func GetAdapter(adapterID string) (*HcitoolDev, error) {

	list, err := GetAdapters()
	if err != nil {
		return nil, err
	}

	for _, a := range list {
		if a.ID == adapterID {
			return a, nil
		}
	}

	return nil, nil
}

// GetAdapters Return a list of adapters using hcitool as backend
func GetAdapters() ([]*HcitoolDev, error) {

	list := make([]*HcitoolDev, 0)
	raw, err := linux.CmdExec("hcitool", "dev")
	if err != nil {
		return nil, err
	}

	// raw:
	// Devices:
	// 	hci1	70:C9:4E:58:AA:7E

	lines := strings.Split(raw, "\n")
	lines = lines[1:]

	// 	hci1	70:C9:4E:58:AA:7E
	re1 := regexp.MustCompile("^[ \t]*([a-zA-Z0-9]+)[ \t]*([a-zA-Z0-9:]+)$")

	for i := 0; i < len(lines); i++ {

		if !re1.MatchString(lines[i]) {
			continue
		}

		el := new(HcitoolDev)

		res := re1.FindStringSubmatch(lines[i])
		if len(res) > 1 {
			el.ID = res[1]
			el.Address = res[2]
			list = append(list, el)
		}

	}

	return list, nil
}

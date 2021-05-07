package hciconfig

import (
	"errors"
	"strings"

	"github.com/muka/go-bluetooth/hw/linux/cmd"
)

// GetAdapters return the list of available adapters
func GetAdapters() ([]HCIConfigResult, error) {

	out, err := cmd.Exec("hciconfig")
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("hciconfig provided no response")
	}

	list := []HCIConfigResult{}
	parts := strings.Split(out, "\nhci")

	for i, el := range parts {
		if i > 0 {
			el = "hci" + el
		}
		cfg := parseControllerInfo(el)
		list = append(list, cfg)
	}

	// log.Debugf("%++v", list)

	return list, nil
}

// GetAdapter return an adapter
func GetAdapter(adapterID string) (*HCIConfigResult, error) {
	h := NewHCIConfig(adapterID)
	return h.Status()
}

// NewHCIConfig initialize a new HCIConfig
func NewHCIConfig(adapterID string) *HCIConfig {
	return &HCIConfig{adapterID}
}

//HCIConfigResult contains details for an adapter
type HCIConfigResult struct {
	AdapterID string
	Enabled   bool
	Address   string
	Type      string
	Bus       string
}

// HCIConfig an hciconfig command wrapper
type HCIConfig struct {
	adapterID string
}

func parseControllerInfo(out string) HCIConfigResult {
	cfg := HCIConfigResult{}

	cfg.AdapterID = strings.Trim(out[:6], " \t:")

	s := strings.Replace(out[6:], "\t", "", -1)
	lines := strings.Split(s, "\n")
	// var parts []string
	for i, line := range lines {
		if i > 2 {
			break
		}
		if i == 2 {
			pp := strings.Split(line, " ")
			cfg.Enabled = (pp[0] == "UP")
			continue
		}

		subparts := strings.Split(line, "  ")
		for _, subpart := range subparts {
			pp := strings.Split(subpart, ": ")
			switch pp[0] {
			case "Type":
				cfg.Type = pp[1]
				continue
			case "Bus":
				cfg.Bus = pp[1]
				continue
			case "BD Address":
				cfg.Address = pp[1]
				continue
			}
		}
	}

	return cfg
}

//Status return status information for a hci device
func (h *HCIConfig) Status() (*HCIConfigResult, error) {

	out, err := cmd.Exec("hciconfig", h.adapterID)
	if err != nil {
		return nil, err
	}

	cfg := parseControllerInfo(out)

	return &cfg, nil
}

// Up Turn on an HCI device
func (h *HCIConfig) Up() (*HCIConfigResult, error) {
	_, err := cmd.Exec("hciconfig", h.adapterID, "up")
	if err != nil {
		return nil, err
	}
	return h.Status()
}

// Down Turn down an HCI device
func (h *HCIConfig) Down() (*HCIConfigResult, error) {
	_, err := cmd.Exec("hciconfig", h.adapterID, "down")
	if err != nil {
		return nil, err
	}
	return h.Status()
}

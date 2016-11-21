package linux

import (
	"bytes"
	"os/exec"
	"strings"
)

func NewHCIConfig(adapterID string) *HCIConfig {
	return &HCIConfig{adapterID}
}

//HCIConfigResult contains details for an adapter
type HCIConfigResult struct {
	Enabled bool
	Address string
	Type    string
	Bus     string
}

// HCIConfig an hciconfig command wrapper
type HCIConfig struct {
	adapterID string
}

//Status return status information for a hci device
func (h *HCIConfig) Status() (HCIConfigResult, error) {

	cfg := HCIConfigResult{}

	cmd := exec.Command("hciconfig", h.adapterID)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return HCIConfigResult{}, err
	}

	s := strings.Replace(out.String()[6:], "\t", "", -1)
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

	return cfg, nil
}

// Up Turn on an HCI device
func (h *HCIConfig) Up() (HCIConfigResult, error) {
	cmd := exec.Command("hciconfig", h.adapterID, "up")
	err := cmd.Run()
	if err != nil {
		return HCIConfigResult{}, err
	}
	return h.Status()
}

// Down Turn down an HCI device
func (h *HCIConfig) Down() (HCIConfigResult, error) {
	cmd := exec.Command("hciconfig", h.adapterID, "down")
	err := cmd.Run()
	if err != nil {
		return HCIConfigResult{}, err
	}
	return h.Status()
}

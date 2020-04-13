package btmgmt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/hw/linux/cmd"
)

const (
	DefaultBinPath = "btmgmt"
)

//BtAdapter contains info about adapter from btmgmt
type BtAdapter struct {
	ID                string
	Name              string
	ShortName         string
	Addr              string
	Version           string
	Manufacturer      string
	Class             string
	SupportedSettings []string
	CurrentSettings   []string
}

//GetAdapter return an adapter
func GetAdapter(adapterID string) (*BtAdapter, error) {
	adapters, err := GetAdapters()
	if err != nil {
		return nil, err
	}
	for _, a := range adapters {
		if a.ID == adapterID {
			return a, nil
		}
	}
	return nil, fmt.Errorf("Adapter %s not found", adapterID)
}

//GetAdapters return a list of adapters
func GetAdapters() ([]*BtAdapter, error) {

	raw, err := cmd.Exec("btmgmt", "info")
	if err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return nil, errors.New("btmgmt provided no response")
	}

	list := make([]*BtAdapter, 0)
	lines := strings.Split(raw, "\n")
	lines = lines[1:]

	// hci1:	Primary controller
	re1 := regexp.MustCompile("([a-z0-9]+):[ ]*")
	//	addr 10:08:B1:72:F5:98 version 6 manufacturer 93 class 0x000000
	re2 := regexp.MustCompile("\taddr ([a-zA-z0-9:]+) version ([0-9]+) manufacturer ([0-9]+) class ([0-9a-zA-Zx]+)")
	//	supported settings:
	re3 := regexp.MustCompile("\t.+: (.*)")
	//	(short )name
	re4 := regexp.MustCompile("\t.*name (.*)")

	// track if parsing an adapter
	for i := 0; i < len(lines); i++ {

		if !re1.MatchString(lines[i]) {
			continue
		}

		el := new(BtAdapter)

		res := re1.FindStringSubmatch(lines[i])
		if len(res) > 1 {
			el.ID = res[1]
		}

		i++
		res = re2.FindStringSubmatch(lines[i])
		if len(res) > 1 {
			el.Addr = res[1]
			el.Version = res[2]
			el.Manufacturer = res[3]
			el.Class = res[4]
		}
		i++
		//supported settings
		res = re3.FindStringSubmatch(lines[i])
		if len(res) > 1 {
			el.SupportedSettings = strings.Split(res[1], " ")
		}

		i++
		//current settings
		res = re3.FindStringSubmatch(lines[i])
		if len(res) > 1 {
			el.CurrentSettings = strings.Split(res[1], " ")
		}

		if lines[i] == "" {
			//empty line
			i++
		}

		i++
		//name
		res = re4.FindStringSubmatch(lines[i])
		if len(res) > 1 {
			el.Name = res[1]
		}

		i++
		//short name
		res = re4.FindStringSubmatch(lines[i])
		if len(res) > 1 {
			el.ShortName = res[1]
		}

		list = append(list, el)
	}

	return list, nil
}

// NewBtMgmt init a new BtMgmt command
func NewBtMgmt(adapterID string) *BtMgmt {
	return &BtMgmt{adapterID, DefaultBinPath}
}

// BtMgmt btmgmt command wrapper
type BtMgmt struct {
	adapterID string
	// BinPath configure the CLI path to btmgmt
	BinPath string
}

// btmgmt cmd wrapper
func (h *BtMgmt) cmd(args ...string) error {
	cmdArgs := []string{h.BinPath, "--index", h.adapterID}
	cmdArgs = append(cmdArgs, args...)
	_, err := cmd.Exec(cmdArgs...)
	if err != nil {
		return err
	}
	return nil
}

// Reset reset the power
func (h *BtMgmt) Reset() error {
	err := h.SetPowered(false)
	if err != nil {
		return err
	}
	return h.SetPowered(true)
}

// SetDeviceID Set Device ID name
func (h *BtMgmt) SetDeviceID(did string) error {
	return h.cmd("did", did)
}

// SetName Set local name
func (h *BtMgmt) SetName(name string) error {
	return h.cmd("name", name)
}

// SetClass set device class
func (h *BtMgmt) SetClass(major, minor string) error {
	return h.cmd("class", major, minor)
}

// SetPowered set power to adapter
func (h *BtMgmt) setFlag(flag string, val bool) error {
	var v string
	if val {
		v = "on"
	} else {
		v = "off"
	}
	return h.cmd(flag, v)
}

// SetPowered set power to adapter
func (h *BtMgmt) SetPowered(status bool) error {
	return h.setFlag("power", status)
}

// SetDiscoverable  Set discoverable state
func (h *BtMgmt) SetDiscoverable(status bool) error {
	return h.setFlag("discov", status)
}

// SetConnectable    	Set connectable state
func (h *BtMgmt) SetConnectable(status bool) error {
	return h.setFlag("connectable", status)
}

// SetFastConnectable 	Set fast connectable state
func (h *BtMgmt) SetFastConnectable(status bool) error {
	return h.setFlag("fast", status)
}

// SetBondable  	Set bondable state
func (h *BtMgmt) SetBondable(status bool) error {
	return h.setFlag("bondable", status)
}

// SetPairable  	Set bondable state
func (h *BtMgmt) SetPairable(status bool) error {
	return h.setFlag("pairable", status)
}

// SetLinkLevelSecurity Set link level security
func (h *BtMgmt) SetLinkLevelSecurity(status bool) error {
	return h.setFlag("linksec", status)
}

// SetSsp     Set SSP mode
func (h *BtMgmt) SetSsp(status bool) error {
	return h.setFlag("ssp", status)
}

// SetSc Toogle SC support
func (h *BtMgmt) SetSc(status bool) error {
	return h.setFlag("sc", status)
}

// SetHs Set HS support
func (h *BtMgmt) SetHs(status bool) error {
	return h.setFlag("hs", status)
}

// SetLe Set LE support
func (h *BtMgmt) SetLe(status bool) error {
	return h.setFlag("le", status)
}

// SetAdvertising    	Set LE advertising
func (h *BtMgmt) SetAdvertising(status bool) error {
	return h.setFlag("advertising", status)
}

// SetBredr   Set BR/EDR support
func (h *BtMgmt) SetBredr(status bool) error {
	return h.setFlag("bredr", status)
}

// SetPrivacy Set privacy support
func (h *BtMgmt) SetPrivacy(status bool) error {
	return h.setFlag("privacy", status)
}

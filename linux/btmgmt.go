package linux

import (
	"regexp"
	"strings"
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

//GetAdapters return a list of adapters
func GetAdapters() ([]BtAdapter, error) {

	raw, err := CmdExec("btmgmt", "info")
	if err != nil {
		return nil, err
	}

	list := []BtAdapter{}
	lines := strings.Split(raw, "\n")
	lines = lines[1:]

	// hci1:	Primary controller
	re1 := regexp.MustCompile("([a-z0-1]+):[ ]*")
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

		el := BtAdapter{}

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
	return &BtMgmt{adapterID}
}

// BtMgmt an hciconfig command wrapper
type BtMgmt struct {
	adapterID string
}

func setFlag(idx string, flag string, val bool) error {
	var v string
	if val {
		v = "on"
	} else {
		v = "off"
	}
	_, err := CmdExec("btmgmt", "--index", idx, flag, v)
	if err != nil {
		return err
	}
	return nil
}

// SetDeviceID Set Device ID name
func (h *BtMgmt) SetDeviceID(did string) error {
	_, err := CmdExec("btmgmt", "--index", h.adapterID, "did", did)
	if err != nil {
		return err
	}
	return nil
}

// SetName Set local name
func (h *BtMgmt) SetName(name string) error {
	_, err := CmdExec("btmgmt", "--index", h.adapterID, "name", name)
	if err != nil {
		return err
	}
	return nil
}

// SetClass set device class
func (h *BtMgmt) SetClass(major, minor string) error {
	_, err := CmdExec("btmgmt", "--index", h.adapterID, "class", major, minor)
	if err != nil {
		return err
	}
	return nil
}

// TogglePowered set power to adapter
func (h *BtMgmt) TogglePowered(status bool) error {
	return setFlag("power", h.adapterID, status)
}

// ToggleDiscoverable  Toggle discoverable state
func (h *BtMgmt) ToggleDiscoverable(status bool) error {
	return setFlag("discov", h.adapterID, status)
}

// ToggleConnectable    	Toggle connectable state
func (h *BtMgmt) ToggleConnectable(status bool) error {
	return setFlag("connectable", h.adapterID, status)
}

// ToggleFastConnectable 	Toggle fast connectable state
func (h *BtMgmt) ToggleFastConnectable(status bool) error {
	return setFlag("fast", h.adapterID, status)
}

// ToggleBondable  	Toggle bondable state
func (h *BtMgmt) ToggleBondable(status bool) error {
	return setFlag("bondable", h.adapterID, status)
}

// TogglePairable  	Toggle bondable state
func (h *BtMgmt) TogglePairable(status bool) error {
	return setFlag("pairable", h.adapterID, status)
}

// ToggleLinkLevelSecurity Toggle link level security
func (h *BtMgmt) ToggleLinkLevelSecurity(status bool) error {
	return setFlag("linksec", h.adapterID, status)
}

// ToggleSsp     Toggle SSP mode
func (h *BtMgmt) ToggleSsp(status bool) error {
	return setFlag("ssp", h.adapterID, status)
}

// ToggleSc Toogle SC support
func (h *BtMgmt) ToggleSc(status bool) error {
	return setFlag("sc", h.adapterID, status)
}

// ToggleHs Toggle HS support
func (h *BtMgmt) ToggleHs(status bool) error {
	return setFlag("hs", h.adapterID, status)
}

// ToggleLe Toggle LE support
func (h *BtMgmt) ToggleLe(status bool) error {
	return setFlag("le", h.adapterID, status)
}

// ToggleAdvertising    	Toggle LE advertising
func (h *BtMgmt) ToggleAdvertising(status bool) error {
	return setFlag("advertising", h.adapterID, status)
}

// ToggleBredr   Toggle BR/EDR support
func (h *BtMgmt) ToggleBredr(status bool) error {
	return setFlag("bredr", h.adapterID, status)
}

// TogglePrivacy Toggle privacy support
func (h *BtMgmt) TogglePrivacy(status bool) error {
	return setFlag("privacy", h.adapterID, status)
}

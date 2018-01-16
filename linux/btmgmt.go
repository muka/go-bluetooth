package linux

// NewBtMgmt init a new BtMgmt command
func NewBtMgmt(adapterID string) *BtMgmt {
	return &BtMgmt{adapterID}
}

//BtMgmtResult contains details for an adapter
type BtMgmtResult struct {
	Enabled bool
	Address string
	Type    string
	Bus     string
}

// BtMgmt an hciconfig command wrapper
type BtMgmt struct {
	adapterID string
}

func setFlag(flag string, val bool) error {
	var v string
	if val {
		v = "on"
	} else {
		v = "off"
	}
	_, err := CmdExec("btmgmt", flag, v)
	if err != nil {
		return err
	}
	return nil
}

// SetDeviceID Set Device ID name
func (h *BtMgmt) SetDeviceID(did string) error {
	_, err := CmdExec("did", did)
	if err != nil {
		return err
	}
	return nil
}

// SetName Set local name
func (h *BtMgmt) SetName(name string) error {
	_, err := CmdExec("name", name)
	if err != nil {
		return err
	}
	return nil
}

// SetClass set device class
func (h *BtMgmt) SetClass(major, minor string) error {
	_, err := CmdExec("class", major, minor)
	if err != nil {
		return err
	}
	return nil
}

// TogglePowered set power to adapter
func (h *BtMgmt) TogglePowered(status bool) error {
	return setFlag("power", status)
}

// ToggleDiscoverable  Toggle discoverable state
func (h *BtMgmt) ToggleDiscoverable(status bool) error {
	return setFlag("discov", status)
}

// ToggleConnectable    	Toggle connectable state
func (h *BtMgmt) ToggleConnectable(status bool) error {
	return setFlag("connectable", status)
}

// ToggleFastConnectable 	Toggle fast connectable state
func (h *BtMgmt) ToggleFastConnectable(status bool) error {
	return setFlag("fast", status)
}

// ToggleBondable  	Toggle bondable state
func (h *BtMgmt) ToggleBondable(status bool) error {
	return setFlag("bondable", status)
}

// TogglePairable  	Toggle bondable state
func (h *BtMgmt) TogglePairable(status bool) error {
	return setFlag("pairable", status)
}

// ToggleLinkLevelSecurity Toggle link level security
func (h *BtMgmt) ToggleLinkLevelSecurity(status bool) error {
	return setFlag("linksec", status)
}

// ToggleSsp     Toggle SSP mode
func (h *BtMgmt) ToggleSsp(status bool) error {
	return setFlag("ssp", status)
}

// ToggleSc Toogle SC support
func (h *BtMgmt) ToggleSc(status bool) error {
	return setFlag("sc", status)
}

// ToggleHs Toggle HS support
func (h *BtMgmt) ToggleHs(status bool) error {
	return setFlag("hs", status)
}

// ToggleLe Toggle LE support
func (h *BtMgmt) ToggleLe(status bool) error {
	return setFlag("le", status)
}

// ToggleAdvertising    	Toggle LE advertising
func (h *BtMgmt) ToggleAdvertising(status bool) error {
	return setFlag("advertising", status)
}

// ToggleBredr   Toggle BR/EDR support
func (h *BtMgmt) ToggleBredr(status bool) error {
	return setFlag("bredr", status)
}

// TogglePrivacy Toggle privacy support
func (h *BtMgmt) TogglePrivacy(status bool) error {
	return setFlag("privacy", status)
}

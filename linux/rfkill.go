package linux

// Code based on rfkill.go from skycoin/skycoin project
// https://github.com/skycoin/skycoin/blob/master/src/aether/wifi/linux/rfkill.go

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func limitText(text []byte) string {
	t := strings.TrimSpace(string(text))
	if len(t) > 150 {
		t = t[0:150] + "..."
	}
	return "[" + t + "]"
}

// RFKill is a wrapper for linux utility: rfkill
// Checks the status of kill switches. If either is set, the device will be disabled.
// Soft = Software (set by software)
// Hard = Hardware (physical on/off switch on the device)
// Identifiers = all, wifi, wlan, bluetooth, uwb, ultrawideband, wimax, wwan, gps, fm
// See: http://wireless.kernel.org/en/users/Documentation/rfkill
type RFKill struct{}

// NewRFKill Creates a new RFKill instance
func NewRFKill() RFKill {
	return RFKill{}
}

// RFKillResult Result of rfkill request
type RFKillResult struct {
	Index          int
	IdentifierType string
	Description    string
	SoftBlocked    bool
	HardBlocked    bool
}

// IsInstalled Checks if the program rfkill exists using PATH environment variable
func (rfkill RFKill) IsInstalled() bool {
	_, err := exec.LookPath("rfkill")
	if err != nil {
		return false
	}
	return true
}

// ListAll Returns a list of rfkill results for every identifier type
func (rfkill RFKill) ListAll() ([]RFKillResult, error) {
	rfks := []RFKillResult{}
	rfk := RFKillResult{}
	fq := rfkill.fileQuery

	// instead of parsing "rfkill list", query the filesystem
	dirInfos, err := ioutil.ReadDir("/sys/class/rfkill/")
	if err != nil {
		return nil, fmt.Errorf(
			"RFKill: Error reading directory '/sys/class/rfkill/': %v", err)
	}

	for _, dirInfo := range dirInfos {
		// directory starts with "rfkill"
		if len(dirInfo.Name()) > 6 && dirInfo.Name()[0:6] == "rfkill" {
			qp := "/sys/class/rfkill/" + dirInfo.Name()

			rfk.Index, _ = strconv.Atoi(fq(qp + "/index"))
			rfk.IdentifierType = fq(qp + "/type")
			rfk.Description = fq(qp + "/name")
			rfk.SoftBlocked = false
			rfk.HardBlocked = false
			if fq(qp+"/soft") == "1" {
				rfk.SoftBlocked = true
			}
			if fq(qp+"/hard") == "1" {
				rfk.HardBlocked = true
			}

			rfks = append(rfks, rfk)
			rfk = RFKillResult{}
		}
	}
	return rfks, nil
}

// SoftBlock RFKill Sets a software block on an identifier
func (rfkill RFKill) SoftBlock(identifier string) error {

	cmd := exec.Command("rfkill", "block", identifier)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Command Error: %v : %v", err, limitText(out))
		return err
	}

	return nil
}

//SoftUnblock Removes a software block on an identifier
func (rfkill RFKill) SoftUnblock(identifier string) error {

	cmd := exec.Command("rfkill", "unblock", identifier)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Command Error: %v : %v", err, limitText(out))
		return err
	}

	return nil
}

//IsBlocked Checks if an identifier has a software or hardware block
func (rfkill RFKill) IsBlocked(identifier string) bool {
	rfks, _ := rfkill.ListAll()
	for _, rfk := range rfks {
		if rfkill.checkThis(rfk, identifier) {
			if rfk.SoftBlocked || rfk.HardBlocked {
				return true
			}
		}
	}
	return false
}

//IsSoftBlocked Checks if an identifier has a software block
func (rfkill RFKill) IsSoftBlocked(identifier string) bool {
	rfks, _ := rfkill.ListAll()
	for _, rfk := range rfks {
		if rfkill.checkThis(rfk, identifier) {
			if rfk.SoftBlocked {
				return true
			}
		}
	}
	return false
}

//IsHardBlocked Checks if an identifier has a hardware block
func (rfkill RFKill) IsHardBlocked(identifier string) bool {
	rfks, _ := rfkill.ListAll()
	for _, rfk := range rfks {
		if rfkill.checkThis(rfk, identifier) {
			if rfk.HardBlocked {
				return true
			}
		}
	}
	return false
}

//IsBlockedAfterUnblocking Checks if an identifier has a software or hardware block after
// removing a software block if it exists
func (rfkill RFKill) IsBlockedAfterUnblocking(identifier string) bool {
	if rfkill.IsBlocked(identifier) {
		rfkill.SoftUnblock(identifier)
		if rfkill.IsBlocked(identifier) {
			return true
		}
	}

	return false
}

func (rfkill RFKill) checkThis(rfk RFKillResult, identifier string) bool {
	switch identifier {
	case "":
		return true
	case "all":
		return true
	case rfk.IdentifierType:
		return true
	}
	return false
}

func (rfkill RFKill) fileQuery(queryFile string) string {
	out, _ := ioutil.ReadFile(queryFile)
	outs := string(out)
	outs = strings.TrimSpace(outs)
	return outs
}

package linux

import (
	"bytes"
	"errors"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// CmdExec Execute a command
func CmdExec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Debugf("Exec: %v", args)

	cmd := exec.Command(baseCmd, cmdArgs...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		out := errbuf.String()
		if out != "" {
			return "", errors.New(out)
		}
		return "", err
	}

	return outbuf.String(), nil
}

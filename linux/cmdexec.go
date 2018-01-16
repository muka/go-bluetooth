package linux

import (
	"bytes"
	"errors"
	"os/exec"
)

// CmdExec Execute a command
func CmdExec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		out := errbuf.String()
		err = errors.New(string(out))
		return "", err
	}

	return outbuf.String(), nil
}

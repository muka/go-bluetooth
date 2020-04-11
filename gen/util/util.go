package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Mkdir Create a dir if not exists
func Mkdir(dirpath string) error {
	err := os.Mkdir(dirpath, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// ListFiles return a list of bluez api txt
func ListFiles(dir string) ([]string, error) {

	list := make([]string, 0)

	if !Exists(dir) {
		return list, fmt.Errorf("Doc dir not found %s", dir)
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, "mgmt-api.txt") {
			return nil
		}

		if !strings.HasSuffix(path, "-api.txt") {
			return nil
		}

		list = append(list, path)
		return nil
	})

	if err != nil {
		log.Errorf("Failed to list files: %s", err)
		return list, nil
	}

	return list, nil
}

// ReadFile read a file content
func ReadFile(srcFile string) ([]byte, error) {
	file, err := os.Open(srcFile)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//GetGitVersion return the docs git version
func GetGitVersion(docsDir string) (string, error) {
	cmd := exec.Command("git", "describe")
	cmd.Dir = docsDir
	res, err := cmd.CombinedOutput()
	return strings.Trim(string(res), " \n\r"), err
}

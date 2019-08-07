package gen

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Parse bluez DBus API docs
func Parse(docsDir string, filters []string) (BluezAPI, error) {
	files, err := ListFiles(docsDir)
	if err != nil {
		return BluezAPI{}, err
	}
	apis := []ApiGroup{}
	for _, file := range files {

		keep := true
		if len(filters) > 0 {
			keep = false
			for _, filter := range filters {
				if strings.Contains(file, filter) {
					keep = true
					break
				}
			}
		}

		if !keep {
			continue
		}

		apiGroup, err := NewApiGroup(file)
		if err != nil {
			log.Errorf("Failed to load %s, skipped", file)
			continue
		}
		apis = append(apis, apiGroup)
	}

	version, err := getGitVersion(docsDir)
	if err != nil {
		log.Errorf("Failed to parse version: %s", err)
	}

	return BluezAPI{
		Version: version,
		Api:     apis,
	}, nil
}

func getGitVersion(docsDir string) (string, error) {
	cmd := exec.Command("git", "describe")
	cmd.Dir = docsDir
	res, err := cmd.CombinedOutput()
	return strings.Trim(string(res), " \n\r"), err
}

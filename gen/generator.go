package gen

import (
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Generate(apiGroups []ApiGroup, outDir string) error {

	err := mkdir(outDir)
	if err != nil {
		log.Errorf("Failed to mkdir %s: %s", outDir, err)
		return err
	}

	filename := filepath.Join(outDir, "errors.go")
	err = ErrorsTemplate(filename, apiGroups)
	if err != nil {
		return err
	}

	for _, apiGroup := range apiGroups {

		apiName := strings.Replace(apiGroup.FileName, "-api.txt", "", -1)
		// log.Debugf("--- Generating %s API ---", apiName)

		dirpath := path.Join(outDir, apiName)
		err := mkdir(dirpath)
		if err != nil {
			log.Errorf("Failed to mkdir %s: %s", dirpath, err)
			continue
		}

		// RootTemplate(apiGroup)

	}

	return nil
}

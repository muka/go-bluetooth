package generator

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/muka/go-bluetooth/gen"
	log "github.com/sirupsen/logrus"
)

func Generate(apiGroups []gen.ApiGroup, outDir string) error {

	err := gen.Mkdir(outDir)
	if err != nil {
		log.Errorf("Failed to mkdir %s: %s", outDir, err)
		return err
	}

	outDir += "/profile"
	err = gen.Mkdir(outDir)
	if err != nil {
		log.Errorf("Failed to mkdir %s: %s", outDir, err)
		return err
	}

	filename := filepath.Join(outDir, "errors.go")
	err = ErrorsTemplate(filename, apiGroups)
	if err != nil {
		return err
	}

	// filename = filepath.Join(outDir, "interfaces.go")
	// err = InterfacesTemplate(filename, apiGroups)
	// if err != nil {
	// 	return err
	// }

	for _, apiGroup := range apiGroups {

		apiName := getApiPackage(apiGroup)
		dirpath := path.Join(outDir, apiName)
		err := gen.Mkdir(dirpath)
		if err != nil {
			log.Errorf("Failed to mkdir %s: %s", dirpath, err)
			continue
		}

		rootFile := path.Join(dirpath, "gen_"+apiName+".go")
		if !gen.Exists(rootFile) {
			err = RootTemplate(rootFile, apiGroup)
			if err != nil {
				log.Errorf("Failed to create %s: %s", rootFile, err)
				continue
			}
			log.Debugf("Wrote %s", rootFile)
		} else {
			log.Infof("Skipped, file exists: %s", rootFile)
		}

		for _, api := range apiGroup.Api {

			pts := strings.Split(api.Interface, ".")
			apiBaseName := pts[len(pts)-1]

			apiFilename := path.Join(dirpath, fmt.Sprintf("%s.go", apiBaseName))
			apiGenFilename := path.Join(dirpath, fmt.Sprintf("gen_%s.go", apiBaseName))

			if gen.Exists(apiFilename) {
				log.Infof("Skipped generation, API file exists: %s", apiFilename)
				continue
			}

			if gen.Exists(apiGenFilename) {
				log.Infof("Skipped, file exists: %s", apiGenFilename)
				continue
			}

			err1 := ApiTemplate(apiGenFilename, api, apiGroup)
			if err1 != nil {
				log.Errorf("Api generation failed %s: %s", api.Title, err1)
				return err1
			}
			log.Debugf("Wrote %s", apiGenFilename)
		}

	}

	return nil
}

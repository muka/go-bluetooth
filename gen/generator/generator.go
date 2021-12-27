package generator

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/muka/go-bluetooth/gen"
	"github.com/muka/go-bluetooth/gen/util"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/imports"
)

// Generate go code from the API definition
func Generate(bluezApi gen.BluezAPI, outDir string, debug bool, forceOverwrite bool) error {

	apiGroups := bluezApi.Api

	err := util.Mkdir(outDir)
	if err != nil {
		log.Errorf("Failed to mkdir %s: %s", outDir, err)
		return err
	}

	outDir += "/profile"
	err = util.Mkdir(outDir)
	if err != nil {
		log.Errorf("Failed to mkdir %s: %s", outDir, err)
		return err
	}

	errorsFile := path.Join(outDir, "gen_errors.go")
	if forceOverwrite || !util.Exists(errorsFile) {
		err = ErrorsTemplate(errorsFile, apiGroups)
		if err != nil {
			return err
		}
	}

	versionFile := path.Join(outDir, "gen_version.go")
	if forceOverwrite || !util.Exists(versionFile) {
		err = VersionTemplate(versionFile, bluezApi.Version)
		if err != nil {
			return err
		}
	}

	// filename = filepath.Join(outDir, "interfaces.go")
	// err = InterfacesTemplate(filename, apiGroups)
	// if err != nil {
	// 	return err
	// }

	for _, apiGroup := range apiGroups {

		if apiGroup == nil {
			continue
		}

		apiName := getApiPackage(apiGroup)
		dirpath := path.Join(outDir, apiName)
		err := util.Mkdir(dirpath)
		if err != nil {
			log.Errorf("Failed to mkdir %s: %s", dirpath, err)
			continue
		}

		rootFile := path.Join(dirpath, fmt.Sprintf("gen_%s.go", apiName))

		if forceOverwrite || !util.Exists(rootFile) {
			err = RootTemplate(rootFile, apiGroup)
			if err != nil {
				log.Errorf("Failed to create %s: %s", rootFile, err)
				continue
			}
			if debug {
				log.Tracef("Wrote %s", rootFile)
			}
		}

		for _, api := range apiGroup.Api {

			if api == nil {
				continue
			}

			pts := strings.Split(api.Interface, ".")
			apiBaseName := pts[len(pts)-1]
			apiBaseName = strings.Replace(apiBaseName, " [experimental]", "", -1)

			apiFilename := path.Join(dirpath, fmt.Sprintf("%s.go", apiBaseName))
			apiGenFilename := path.Join(dirpath, fmt.Sprintf("gen_%s.go", apiBaseName))

			if util.Exists(apiFilename) {
				// log.Debugf("Skipped generation, API file exists: %s", apiFilename)
				continue
			}

			if !forceOverwrite && util.Exists(apiGenFilename) {
				// log.Debugf("Skipped, file exists: %s", apiGenFilename)
				continue
			}

			err1 := ApiTemplate(apiGenFilename, api, apiGroup)
			if err1 != nil {
				log.Errorf("Api generation failed %s: %s", api.Title, err1)
				return err1
			}
			if debug {
				log.Tracef("Wrote %s", apiGenFilename)
			}

			code, err := imports.Process(apiGenFilename, nil, nil)
			if err != nil {
				log.Tracef("format code: %s: %v", apiGenFilename, err)
			}

			if err := ioutil.WriteFile(apiGenFilename, code, 0644); err != nil {
				log.Tracef("rewrite with formatted code: %s", apiGenFilename)
				return err
			}
		}
	}
	return nil
}

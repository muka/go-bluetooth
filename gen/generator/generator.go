package generator

import (
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

		rootFile := path.Join(dirpath, apiName+".go")
		err = RootTemplate(rootFile, apiGroup)
		if err != nil {
			log.Errorf("Failed to create %s: %s", rootFile, err)
			continue
		}

		for _, api := range apiGroup.Api {

			pts := strings.Split(api.Interface, ".")
			apiFilename := path.Join(dirpath, pts[len(pts)-1]+".go")

			err1 := ApiTemplate(apiFilename, api, apiGroup)
			if err1 != nil {
				log.Errorf("Api generation failed %s: %s", api.Title, err1)
				return err1
			}
		}

	}

	return nil
}

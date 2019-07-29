package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/muka/go-bluetooth/gen"
)

func RootTemplate(filename string, api gen.ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	apidoc := gen.ApiGroupDoc{
		ApiGroup: api,
		Package:  getApiPackage(api),
	}

	apidoc.ApiGroup.Description = prepareDocs(apidoc.ApiGroup.Description, false, 0)

	tmpl := loadtpl("root")
	err = tmpl.Execute(fw, apidoc)
	if err != nil {
		return fmt.Errorf("write tpl: %s", err)
	}

	// log.Debugf("Created %s", filename)
	return nil
}

func ErrorsTemplate(filename string, apis []gen.ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	errors := []string{}
	for _, apiGroup := range apis {
		for _, api := range apiGroup.Api {
			for _, method := range api.Methods {
				for _, err := range method.Errors {
					errors = appendIfMissing(errors, err)
				}
			}
		}
	}

	errorsList := gen.BluezErrors{
		List: make([]gen.BluezError, len(errors)),
	}

	for i, err := range errors {
		errorsList.List[i] = gen.BluezError{
			Name: strings.Replace(err, "org.bluez.Error.", "", 1),
		}
	}

	tmpl := loadtpl("errors")
	err = tmpl.Execute(fw, errorsList)
	if err != nil {
		return fmt.Errorf("tpl write: %s", err)
	}

	// log.Debugf("Created %s", filename)
	return nil
}

func InterfacesTemplate(filename string, apis []gen.ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	interfaces := []gen.InterfaceDoc{}
	for _, apiGroup := range apis {
		for _, api := range apiGroup.Api {

			pts := strings.Split(api.Interface, ".")
			ifaceName := pts[len(pts)-1]
			// org.bluez.obex.AgentManager1
			if len(pts) > 3 {
				ifaceName = ""
				for _, pt := range pts[2:] {
					ifaceName += strings.ToUpper(string(pt[0])) + pt[1:]
				}
			}

			iface := gen.InterfaceDoc{
				Title:     api.Title,
				Name:      ifaceName,
				Interface: api.Interface,
			}
			interfaces = append(interfaces, iface)
		}
	}

	ifaces := gen.InterfacesDoc{
		Interfaces: interfaces,
	}

	tmpl := loadtpl("interfaces")
	err = tmpl.Execute(fw, ifaces)
	if err != nil {
		return fmt.Errorf("tpl write: %s", err)
	}

	// log.Debugf("Created %s", filename)
	return nil
}

package gen

import (
	"fmt"
	"os"
	"strings"

	"github.com/muka/go-bluetooth/gen/override"
	log "github.com/sirupsen/logrus"
)

func ApiTemplate(filename string, api Api, apiGroup ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	apiName := getApiPackage(apiGroup)

	imports := []string{
		"sync",
		"github.com/muka/go-bluetooth/bluez",
	}

	// Expose Properties interface ?
	exposeProps := override.ExposeProperties(api.Interface)

	if exposeProps {
		imports = append(imports, "github.com/fatih/structs")
	}

	// flag to import dbus
	// NOTE: set to true to handle dbus.Signalling
	// importDbus := false

	importDbus := true

	pts := strings.Split(api.Interface, ".")
	iface := pts[len(pts)-1]

	props := []PropertyDoc{}
	for _, p := range api.Properties {

		prop := PropertyDoc{
			Property: p,
		}

		prop.Name = strings.Trim(p.Name, ": \t")
		prop.Property.Docs = prepareDocs(p.Docs, true, 2)
		prop.Property.Type = castType(p.Type)

		val, found := override.GetPropertyTypeOverride(api.Interface, prop.Name)
		if found {
			prop.Property.Type = val
		}

		if !importDbus {
			importDbus = strings.Contains(prop.Property.Type, "dbus.")
		}

		props = append(props, prop)
	}

	methods := []MethodDoc{}
	for _, m := range api.Methods {

		args := []string{}
		params := []string{}
		for _, a := range m.Args {
			arg := a.Name + " " + castType(a.Type)
			args = append(args, arg)
			params = append(params, a.Name)
		}

		mm := MethodDoc{
			Method:     m,
			ArgsList:   strings.Join(args, ", "),
			ParamsList: strings.Join(params, ", "),
		}

		if !importDbus {
			importDbus = strings.Contains(mm.ArgsList, "dbus.")
			log.Debugf("%t %s", importDbus, mm.ArgsList)
		}
		if !importDbus {
			importDbus = strings.Contains(mm.ParamsList, "dbus.")
		}

		mm.Method.Name = strings.Replace(mm.Method.Name, " (optional)", "", -1)
		mm.Method.Docs = prepareDocs(mm.Method.Docs, true, 0)
		mm.Method.ReturnType = castType(mm.Method.ReturnType)

		mm.SingleReturn = len(mm.Method.ReturnType) == 0

		if mm.SingleReturn {
			mm.Method.ReturnType = "error"
		} else {

			// log.Debugf("With return type %s", mm.Method.ReturnType)

			returnTypes := strings.Split(mm.Method.ReturnType, ", ")
			defs := []string{}
			refs := []string{}
			list := []string{}
			for i, returnType := range returnTypes {

				// objInitialization1 := ""
				// objInitialization2 := ""
				// if strings.Contains(returnType, "[]") {
				// 	objInitialization1 = "="
				// 	objInitialization2 = "{}"
				// }

				varName := fmt.Sprintf("val%d", i)
				// def := fmt.Sprintf("var %s %s %s%s", varName, objInitialization1, returnType, objInitialization2)
				def := fmt.Sprintf("var %s %s", varName, returnType)
				ref := "&" + varName

				defs = append(defs, def)
				refs = append(refs, ref)
				list = append(list, varName)

			}

			mm.ReturnVarsDefinition = strings.Join(defs, "\n  ")
			mm.ReturnVarsRefs = strings.Join(refs, ", ")
			mm.ReturnVarsList = strings.Join(list, ", ")

			if !importDbus {
				importDbus = strings.Contains(mm.Method.ReturnType, "dbus.")
			}

			mm.Method.ReturnType = "(" + mm.Method.ReturnType + ", error)"
		}

		if len(mm.Method.Name) == 0 {
			continue
		}

		methods = append(methods, mm)
	}

	if importDbus {
		imports = append(imports, "github.com/godbus/dbus")
	}

	importsTpl := ""
	if len(imports) > 0 {
		for i := range imports {
			imports[i] = fmt.Sprintf(`"%s"`, imports[i])
		}
		importsTpl = fmt.Sprintf("import (\n  %s\n)", strings.Join(imports, "\n  "))
	}

	api.Description = prepareDocs(api.Description, false, 0)
	api.Title = strings.Trim(api.Title, "\n \t")

	ctrs := createConstructors(api)

	apidocs := ApiDoc{
		Imports:          importsTpl,
		Package:          apiName,
		Api:              api,
		InterfaceName:    iface,
		Properties:       props,
		Methods:          methods,
		Constructors:     ctrs,
		ExposeProperties: exposeProps,
	}

	tmpl := loadtpl("api")
	err = tmpl.Execute(fw, apidocs)
	if err != nil {
		return fmt.Errorf("api tpl: %s", err)
	}

	// log.Debugf("Created %s", filename)
	return nil
}

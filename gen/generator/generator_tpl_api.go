package generator

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/muka/go-bluetooth/gen/override"
	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

func ApiTemplate(filename string, api *types.Api, apiGroup *types.ApiGroup) error {

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
		propsImports := []string{
			// "log github.com/sirupsen/logrus",
			// "reflect",
			// "github.com/fatih/structs",
			"github.com/muka/go-bluetooth/util",
			"github.com/muka/go-bluetooth/props",
		}
		imports = append(imports, propsImports...)
	}

	// flag to import dbus
	// NOTE: set to true to handle dbus.Signalling
	// importDbus := false

	importDbus := true

	pts := strings.Split(api.Interface, ".")
	iface := pts[len(pts)-1]

	propsList := map[string]*types.PropertyDoc{}

	for _, p := range api.Properties {

		prop := types.PropertyDoc{
			Property: p,
		}

		prop.Name = strings.Trim(p.Name, ": \t")
		prop.Property.Docs = prepareDocs(p.Docs, true, 2)
		prop.Property.Type = castType(p.Type)
		prop.RawType = getRawType(prop.Property.Type)

		prop.RawTypeInitializer, err = getRawTypeInitializer(prop.Property.Type)
		if err != nil {
			log.Errorf("%s %s: %s", api.Interface, prop.Name, err)
		}

		propsList[prop.Name] = &prop
	}

	propertiesOverride, found := override.GetPropertiesOverride(api.Interface)
	if found {
		log.Debugf("Found overrides %s", api.Interface)
		for propName, propType := range propertiesOverride {

			var prop *types.PropertyDoc

			if _, ok := propsList[propName]; ok {

				prop = propsList[propName]
				log.Debugf("\tUsing overridden property %s", propName)

				prop.Property.Type = propType
				rawTypeInitializer, err := getRawTypeInitializer(prop.Property.Type)
				if err != nil {
					log.Errorf("Override %s %s: %s", api.Interface, prop.Name, err)
				}

				prop.RawTypeInitializer = rawTypeInitializer
				prop.RawType = getRawType(prop.Property.Type)
				//log.Debugf("props --> %s %s", propName, propType)
			} else {

				rawTypeInitializer, err := getRawTypeInitializer(propType)
				if err != nil {
					log.Errorf("Override %s %s: %s", api.Interface, prop.Name, err)
				}

				prop = &types.PropertyDoc{
					Property: &types.Property{
						Name: propName,
						Type: propType,
					},
					RawType:            getRawType(propType),
					RawTypeInitializer: rawTypeInitializer,
				}
				propsList[propName] = prop
			}

			if !importDbus {
				importDbus = strings.Contains(prop.Property.Type, "dbus.")
			}

		}
	}

	props := []types.PropertyDoc{}
	for _, prop := range propsList {

		// add propery flags
		for _, flag := range prop.Flags {
			if flag == types.FlagReadOnly {
				prop.ReadOnly = true
			}
			if flag == types.FlagWriteOnly {
				prop.WriteOnly = true
			}
			if flag == types.FlagReadWrite {
				prop.ReadWrite = true
			}
		}

		props = append(props, *prop)
	}

	sort.Slice(props, func(i, j int) bool {
		return props[i].Name < props[j].Name
	})

	methods := []types.MethodDoc{}
	for _, m := range api.Methods {

		args := []string{}
		params := []string{}
		for _, a := range m.Args {
			argName := renameReserved(a.Name)
			arg := fmt.Sprintf("%s %s", argName, castType(a.Type))
			args = append(args, arg)
			params = append(params, argName)
		}

		mm := types.MethodDoc{
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

				varName := fmt.Sprintf("val%d", i)
				varDeclaration := "var"

				// handle array
				if strings.HasPrefix(returnType, "[]") {
					varDeclaration = ""
					returnType = fmt.Sprintf(":= %s{}", returnType)
				}

				// def := fmt.Sprintf("var %s %s %s%s", varName, objInitialization1, returnType, objInitialization2)
				def := fmt.Sprintf("%s %s %s", varDeclaration, varName, returnType)
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
		imports = append(imports, "github.com/godbus/dbus/v5")
	}

	api.Description = prepareDocs(api.Description, false, 0)
	api.Title = strings.Trim(api.Title, "\n \t")

	ctrs := createConstructors(api)

	for _, c := range ctrs {
		importFmt := strings.Contains(c.ObjectPath, "fmt.")
		if !importFmt {
			importFmt = strings.Contains(c.Service, "fmt.")
		}
		if importFmt {
			imports = append(imports, "fmt")
		}
	}

	importsTpl := ""
	if len(imports) > 0 {
		for i := range imports {
			pts := strings.Split(imports[i], " ")
			if len(pts) == 1 {
				pts = append(pts, pts[0])
				pts[0] = ""
			}
			imports[i] = fmt.Sprintf(`%s "%s"`, pts[0], pts[1])
		}
		importsTpl = fmt.Sprintf("import (\n  %s\n)", strings.Join(imports, "\n  "))
	}

	apidocs := types.ApiDoc{
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

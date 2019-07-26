package gen

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type BluezError struct {
	Name  string
	Error string
}

type BluezErrors struct {
	List []BluezError
}

type MethodDoc struct {
	Method
	ArgsList             string
	ParamsList           string
	SingleReturn         bool
	ReturnVarsDefinition string
	ReturnVarsRefs       string
	ReturnVarsList       string
}

type InterfaceDoc struct {
	Title     string
	Name      string
	Interface string
}

type InterfacesDoc struct {
	Interfaces []InterfaceDoc
}

type PropertyDoc struct {
	Property
}

type ApiGroupDoc struct {
	ApiGroup
	Package string
}

type ApiDoc struct {
	Api           Api
	InterfaceName string
	Package       string
	Properties    []PropertyDoc
	Methods       []MethodDoc
	Imports       string
}

func loadtpl(name string) *template.Template {
	return template.Must(template.ParseFiles("gen/tpl/" + name + ".go.tpl"))
}

func prepareDocs(src string, skipFirstComment bool, leftpad int) string {

	lines := strings.Split(src, "\n")
	result := []string{}

	comment := "// "
	prefixLen := leftpad + len(comment)
	fmtt := fmt.Sprintf("%%%ds%%s", prefixLen)

	for _, line := range lines {
		line = strings.Trim(line, " \t\r")
		if len(line) == 0 {
			continue
		}

		result = append(result, fmt.Sprintf(fmtt, comment, line))
	}
	if skipFirstComment && len(result) > 0 && len(result[0]) > 3 {
		result[0] = result[0][prefixLen:]
	}
	return strings.Join(result, "\n")
}

func toType(t string) string {
	switch strings.Trim(t, " \t\r\n") {
	case "boolean":
		return "bool"
	case "uint16_t":
		return "uint16"
	case "uint32_t":
		return "uint32"
	case "uint8_t":
		return "uint8"
	case "dict":
		return "map[string]interface{}"
	// check in media-api
	case "properties":
		return "string"
	case "object":
		return "dbus.ObjectPath"
	case "objects":
		return "dbus.ObjectPath"
	case "fd":
		return "int32"
	case "<unknown>":
		return ""
	case "unknown":
		return ""
	case "void":
		return ""
	}
	return t
}

func listCastType(typedef string) string {
	// handle multiple items eg. byte, uint16
	if strings.Contains(typedef, ",") && typedef[:5] != "array" {
		parts := strings.Split(typedef, ", ")
		defs := make([]string, 0)
		for _, part := range parts {
			subtype := castType(part)
			if len(subtype) > 0 {
				defs = append(defs, subtype)
			}
		}
		typedef = strings.Join(defs, ", ")
	}
	return typedef
}

func castType(rawtype string) string {

	if rawtype == "" {
		return ""
	}

	typedef := listCastType(rawtype)

	//eg. array{string} or array{string, foo}
	re := regexp.MustCompile(`array\{([a-zA-Z0-9, ]+)\}`)
	matches := re.FindSubmatch([]byte(rawtype))
	if len(matches) > 0 {
		subtype := string(matches[1])
		subtype = listCastType(subtype)
		typedef = "[]" + toType(subtype)
	}

	typedef = toType(typedef)
	// log.Debugf("type casting %s -> %s\n", rawtype, typedef)

	return typedef
}

func getApiPackage(apiGroup ApiGroup) string {
	apiName := strings.Replace(apiGroup.FileName, "-api.txt", "", -1)
	apiName = strings.Replace(apiName, "-", "_", -1)
	return apiName
}

func RootTemplate(filename string, api ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	apidoc := ApiGroupDoc{
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

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func ErrorsTemplate(filename string, apis []ApiGroup) error {

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

	errorsList := BluezErrors{
		List: make([]BluezError, len(errors)),
	}

	for i, err := range errors {
		errorsList.List[i] = BluezError{
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

func InterfacesTemplate(filename string, apis []ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	interfaces := []InterfaceDoc{}
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

			iface := InterfaceDoc{
				Title:     api.Title,
				Name:      ifaceName,
				Interface: api.Interface,
			}
			interfaces = append(interfaces, iface)
		}
	}

	ifaces := InterfacesDoc{
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

func ApiTemplate(filename string, api Api, apiGroup ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	apiName := getApiPackage(apiGroup)

	imports := []string{
		"github.com/muka/go-bluetooth/bluez",
	}

	// flag to import dbus
	importDbus := false

	pts := strings.Split(api.Interface, ".")
	iface := pts[len(pts)-1]

	props := []PropertyDoc{}
	for _, p := range api.Properties {

		prop := PropertyDoc{
			Property: p,
		}

		prop.Property.Docs = prepareDocs(p.Docs, true, 2)
		prop.Property.Type = castType(p.Type)

		if !importDbus {
			importDbus = strings.Contains(prop.Property.Type, "ObjectPath")
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
			importDbus = strings.Contains(mm.ArgsList, "ObjectPath")
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

				objInitialization1 := ""
				objInitialization2 := ""
				if strings.Contains(returnType, "[]") {
					objInitialization1 = "="
					objInitialization2 = "{}"
				}

				varName := fmt.Sprintf("val%d", i)
				def := fmt.Sprintf("var %s %s %s%s", varName, objInitialization1, returnType, objInitialization2)
				ref := "&" + varName

				defs = append(defs, def)
				refs = append(refs, ref)
				list = append(list, varName)

			}

			mm.ReturnVarsDefinition = strings.Join(defs, "\n")
			mm.ReturnVarsRefs = strings.Join(refs, ", ")
			mm.ReturnVarsList = strings.Join(list, ", ")

			if !importDbus {
				importDbus = strings.Contains(mm.ReturnVarsDefinition, "ObjectPath")
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

	apidocs := ApiDoc{
		Imports:       importsTpl,
		Package:       apiName,
		Api:           api,
		InterfaceName: iface,
		Properties:    props,
		Methods:       methods,
	}

	tmpl := loadtpl("api")
	err = tmpl.Execute(fw, apidocs)
	if err != nil {
		return fmt.Errorf("api tpl: %s", err)
	}

	// log.Debugf("Created %s", filename)
	return nil
}

package gen

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
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
	ArgsList string
}

type InterfaceDoc struct {
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
	switch t {
	case "bool":
	case "boolean":
		return "bool"
	case "dict":
		return "map[string]interface{}"
	case "object":
		return "interface{}"
	case "<unknown>":
	case "unknown":
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
			subtype := castType(strings.Trim(part, " "))
			subtype = strings.Trim(subtype, " \t")
			if len(subtype) > 0 {
				defs = append(defs, subtype)
			}
		}
		typedef = strings.Join(defs, ", ")
	}
	return typedef
}

func castType(typedef string) string {

	if typedef == "" {
		return ""
	}

	// log.Debugf("\n DBUS TYPE ---- %s", typedef)

	typedef = listCastType(typedef)

	// array{string}
	re := regexp.MustCompile(`array\{([a-zA-Z0-9]+)\}`)
	matches := re.FindSubmatch([]byte(typedef))
	if len(matches) > 0 {
		// log.Debugf(" array of ----------------------------------- %s", matches[1])
		subtype := string(matches[1])
		subtype = listCastType(subtype)
		typedef = "[]" + toType(subtype)
	}

	typedef = toType(typedef)
	// log.Debugf(" GO TYPE ---- %s \n", typedef)

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

	log.Debugf("Created %s", filename)
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

	log.Debugf("Created %s", filename)
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

			iface := InterfaceDoc{
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

	log.Debugf("Created %s", filename)
	return nil
}

func ApiTemplate(filename string, api Api, apiGroup ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	apiName := getApiPackage(apiGroup)

	pts := strings.Split(api.Interface, ".")
	iface := pts[len(pts)-1]

	props := []PropertyDoc{}
	for _, p := range api.Properties {

		p.Docs = prepareDocs(p.Docs, true, 2)
		p.Type = castType(p.Type)

		prop := PropertyDoc{
			Property: p,
		}

		props = append(props, prop)
	}

	methods := []MethodDoc{}
	for _, m := range api.Methods {

		args := []string{}
		for _, a := range m.Args {
			arg := a.Name + " " + castType(a.Type)
			args = append(args, arg)
		}

		m.Docs = prepareDocs(m.Docs, true, 2)

		mm := MethodDoc{
			Method:   m,
			ArgsList: strings.Join(args, ", "),
		}

		mm.Method.ReturnType = castType(mm.Method.ReturnType)

		methods = append(methods, mm)
	}

	api.Description = prepareDocs(api.Description, false, 0)
	api.Title = strings.Trim(api.Title, "\n \t")

	apidocs := ApiDoc{
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

	log.Debugf("Created %s", filename)
	return nil
}

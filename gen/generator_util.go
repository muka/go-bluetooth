package gen

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

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
	case "int16":
		return "int16"
	case "uint16_t":
		return "uint16"
	case "uint32_t":
		return "uint32"
	case "uint8_t":
		return "uint8"
	case "dict":
		return "map[string]dbus.Variant"
		// return "map[string]interface{}"
	// check in media-api
	case "properties":
		return "string"
	case "object":
		return "dbus.ObjectPath"
	case "objects":
		return "dbus.ObjectPath"
	case "fd":
		return "dbus.UnixFD"
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

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

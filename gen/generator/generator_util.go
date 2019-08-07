package generator

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/muka/go-bluetooth/gen"
)

var TplPath = "./gen/generator/tpl/%s.go.tpl"

func loadtpl(name string) *template.Template {
	return template.Must(template.ParseFiles(fmt.Sprintf(TplPath, name)))
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
		// return "map[string]dbus.Variant"
		return "map[string]interface{}"
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

func prepareDocs(src string, skipFirstComment bool, leftpad int) string {
	return src
	// lines := strings.Split(src, "\n")
	// result := []string{}
	// // comment := "// "
	// comment := ""
	// prefixLen := leftpad + len(comment)
	// fmtt := fmt.Sprintf("%%%ds%%s", prefixLen)
	//
	// for _, line := range lines {
	// 	line = strings.Trim(line, " \t\r")
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	//
	// 	result = append(result, fmt.Sprintf(fmtt, comment, line))
	// }
	// if skipFirstComment && len(result) > 0 && len(result[0]) > 3 {
	// 	result[0] = result[0][prefixLen:]
	// }
	// return strings.Join(result, "\n")
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

		//special case 1 (obex): `array{string vcard, string name}`
		pts := strings.Split(subtype, ", ")
		pts2 := []string{}
		for _, pt := range pts {
			pts1 := strings.Split(pt, " ")
			// log.Debug("pts1 ", pts1[0])
			if len(pts1) == 2 {
				pts2 = append(pts2, pts1[0])
			}
		}
		if len(pts2) > 0 {
			subtype = strings.Join(pts2, ", ")
		}

		// TODO this is incomplete as it is not handling the case
		// array{ string, string } ==> []string, string
		typedef = "[]" + toType(subtype)

		// log.Debugf("type casting %s -> %s\n", rawtype, typedef)
	}

	typedef = toType(typedef)
	// log.Debugf("type casting %s -> %s\n", rawtype, typedef)

	return typedef
}

func getApiPackage(apiGroup gen.ApiGroup) string {
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

// getRawType clean tag from type
func getRawType(t string) string {
	if strings.Contains(t, "`") {
		p1 := strings.Trim(strings.Split(t, "`")[0], " ")
		return p1
	}
	return t
}

// getRawTypeInitializer return field initializer
func getRawTypeInitializer(t string) string {
	t = getRawType(t)

	// array
	if len(t) >= 2 && t[:2] == "[]" {
		return t + "{}"
	}
	// map
	if len(t) >= 3 && t[:3] == "map" {
		return t + "{}"
	}
	// int*
	if len(t) >= 3 && t[:3] == "int" {
		return t + "(0)"
	}
	// uint*
	if len(t) >= 4 && t[:4] == "uint" {
		return t + "(0)"
	}
	// float*
	if len(t) >= 5 && t[:5] == "float" {
		return t + "(0.0)"
	}

	switch t {
	case "bool":
		return "false"
	case "string":
		return "\"\""
	case "byte":
		return "byte(0)"
		// return "[]uint8{}"
	case "dbus.ObjectPath":
		return "dbus.ObjectPath(\"\")"
	default:
		panic(fmt.Sprintf("Unknown type: %s", t))
	}
}

package generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/gen/override"
)

func toType(t string) string {
	switch strings.ToLower(strings.Trim(t, " \t\r\n")) {
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
	case "variant":
		return "dbus.Variant"
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
	if strings.Contains(typedef, ", ") && typedef[:5] != "array" {
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

	if mappedType, ok := override.MapType(rawtype); ok {
		return mappedType
	}

	typedef := listCastType(rawtype)

	//eg. array{string} or array{string, foo}
	re := regexp.MustCompile(`array\{([a-zA-Z0-9, ()]+)\}`)
	matches := re.FindSubmatch([]byte(rawtype))
	// log.Warnf("submatch -> %s", matches)
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

	// handle case for named type eg. "uint64 token"
	namedTypeRegex := regexp.MustCompile(`([^ ]+) (.*)`)
	namedTypeRegexMatches := namedTypeRegex.FindStringSubmatch(rawtype)
	if len(namedTypeRegexMatches) > 0 {
		typedef = namedTypeRegexMatches[1]
	}

	// log.Debugf("type casting %s -> %s\n", rawtype, typedef)

	return typedef
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
func getRawTypeInitializer(t string) (string, error) {
	t = getRawType(t)

	var checkType = func(text string, match string) bool {
		minlen := len(match)
		return len(text) >= minlen && strings.ToLower(t[:minlen]) == match
	}

	// array
	if checkType(t, "[]") {
		return t + "{}", nil
	}
	// map
	if checkType(t, "map") {
		return t + "{}", nil
	}
	// int*
	if checkType(t, "int") {
		return t + "(0)", nil
	}
	// uint*
	if checkType(t, "uint") {
		return t + "(0)", nil
	}
	// float*
	if checkType(t, "float") {
		return t + "(0.0)", nil
	}

	switch t {
	case "bool":
		return "false", nil
	case "string":
		return "\"\"", nil
	case "byte":
		return "byte(0)", nil
		// return "[]uint8{}"
	case "dbus.ObjectPath":
		return "dbus.ObjectPath(\"\")", nil
	case "dbus.objectpath":
		return "dbus.ObjectPath(\"\")", nil
	case "Track":
		return "Track{}", nil
	}

	return "", fmt.Errorf("unknown type: %s", t)
}

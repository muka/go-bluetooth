package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/gen/types"
	log "github.com/sirupsen/logrus"
)

// NewMethodParser init a MethodParser
func NewMethodParser(debug bool) MethodParser {
	return MethodParser{
		model: new(types.Method),
		debug: debug,
	}
}

//MethodParser wrap a parsable method
type MethodParser struct {
	model *types.Method
	debug bool
}

//Parse a method text
func (g *MethodParser) Parse(raw []byte) (*types.Method, error) {

	var err error = nil
	method := g.model

	re := regexp.MustCompile(`[\t]{1,}(.*?)(?: |\n\t{2,})?(\w+)\(([^)]*)\) ?(.*?)\n((?s).+)`)
	matches1 := re.FindAllSubmatch(raw, -1)

	for _, matches2 := range matches1 {

		rtype := string(matches2[1])
		if len(rtype) > 7 && rtype[:7] == "Methods" {
			rtype = rtype[7:]
		}

		rtype = strings.Trim(rtype, " \t")

		for _, srtype := range strings.Split(rtype, ",") {
			if len(strings.Split(strings.Trim(srtype, " "), " ")) > 2 {
				// log.Warnf("****** %s | %s", strings.Trim(srtype, " "), strings.Split(strings.Trim(srtype, " "), " "))
				return g.model, fmt.Errorf("Method %s return type contains space: `%s`", method.Name, rtype)
			}
		}

		if len(rtype) > 20 {
			log.Warnf("Return type value is too long? `%s`", rtype)
		}

		method.ReturnType = rtype

		name := string(matches2[2])
		method.Name = strings.Trim(name, " \t")

		args := []types.Arg{}
		if len(matches2[3]) > 0 {

			args1 := string(matches2[3])
			if args1 == "void" {
				continue
			}

			argslist := strings.Split(args1, ",")
			for _, arg := range argslist {
				arg = strings.Trim(arg, " ")
				argsparts := strings.Split(arg, " ")
				if argsparts[0] == "void" {
					continue
				}
				if len(argsparts) < 2 {
					if argsparts[0] == "fd" {
						argsparts = []string{"int32", argsparts[0]}
					} else {
						argsparts = []string{"<unknown>", argsparts[0]}
					}
				}

				argType := strings.Trim(argsparts[0], " \t\n")
				arg := types.Arg{
					Type: argType,
					Name: argsparts[1],
				}
				args = append(args, arg)
			}
		}
		method.Args = args
		method.Docs = string(matches2[5])
	}

	//
	re2 := regexp.MustCompile(`(?s)(org\.bluez\.Error\.\w+)`)
	matches2 := re2.FindAllSubmatch(raw, -1)

	if len(matches2) >= 1 {
		for _, merr := range matches2[0] {
			method.Errors = append(method.Errors, string(merr))
		}
	}

	if method.Name == "" {
		return method, errors.New("Empty method name")
	}

	// if strings.Contains(method.ReturnType, "Handle") {
	// 	fmt.Println(method)
	// 	os.Exit(1)
	// }

	if g.debug {
		log.Debugf("\t - %s", method)
	}

	return method, err
}

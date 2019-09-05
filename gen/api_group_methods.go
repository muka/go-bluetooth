package gen

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (g *ApiGroup) parseMethods(raw []byte) []Method {

	methods := make([]Method, 0)
	slices := make([][]byte, 0)

	re := regexp.MustCompile(`(?s)Methods(.+?)(Properties|Signals)[\t ]`)
	matches1 := re.FindSubmatch(raw)

	if len(matches1) == 0 {
		re = regexp.MustCompile(`(?s)[ \t\n]+(.+)`)
		matches1 = re.FindSubmatch(raw)
		if len(matches1) == 1 {
			matches1 = append(matches1, matches1[0])
		}
		// log.Debugf("matches1 %s", matches1[1:])
	}

	// log.Debugf("matches1 %s", matches1[:1])
	// log.Debugf("%s", matches1)

	for _, methodsRaw := range matches1[1:] {

		re1 := regexp.MustCompile(`[ \t]*?(.*?)? ?([^ ]+)\(([^)]+?)?\) ?(.*)`)
		matches2 := re1.FindAllSubmatchIndex(methodsRaw, -1)
		if len(matches2) == 1 {
			if len(methodsRaw) > 0 {
				slices = append(slices, methodsRaw)
			}
		} else {
			prevPos := -1
			for i := 0; i < len(matches2); i++ {

				if prevPos == -1 {
					prevPos = matches2[i][0]
					continue
				}

				nextPos := matches2[i][0]
				methodRaw := methodsRaw[prevPos:nextPos]
				prevPos = nextPos

				if len(methodRaw) > 0 {
					slices = append(slices, methodRaw)
				}

				// keep the last one
				lastItem := len(matches2) - 1
				if i == lastItem {
					methodRaw = methodsRaw[matches2[lastItem][0]:]
					if len(methodRaw) > 0 {
						slices = append(slices, methodRaw)
					}
				}
			}
		}
	}

	if g.debug {
		log.Debug("\tMethods:")
	}
	for _, methodRaw := range slices {
		method, err := g.parseMethod(methodRaw)
		if err != nil {
			log.Warnf("Skip method: %s", err)
			continue
		}
		methods = append(methods, method)
	}

	return methods
}

func (g *ApiGroup) parseMethod(raw []byte) (method Method, err error) {

	re := regexp.MustCompile(`[ \t]*(.*?) ?(\w+)\(([^)]*)\) ?(.*?)\n((?s).+)`)
	matches1 := re.FindAllSubmatch(raw, -1)

	// log.Debugf("matches1 %s", matches1)
	for _, matches2 := range matches1 {

		rtype := string(matches2[1])
		if len(rtype) > 7 && rtype[:7] == "Methods" {
			rtype = rtype[7:]
		}

		rtype = strings.Trim(rtype, " \t")
		for _, srtype := range strings.Split(rtype, ",") {
			if len(strings.Split(strings.Trim(srtype, " "), " ")) > 2 {
				// log.Warnf("****** %s | %s", strings.Trim(srtype, " "), strings.Split(strings.Trim(srtype, " "), " "))
				return method, fmt.Errorf("Method %s return type contains space: `%s`", method.Name, rtype)
			}
		}

		if len(rtype) > 20 {
			log.Warnf("Return type value is too long? `%s`", rtype)
		}

		method.ReturnType = rtype

		name := string(matches2[2])
		method.Name = strings.Trim(name, " \t")

		args := []Arg{}
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
				arg := Arg{
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

	if method.Name != "" && g.debug {
		log.Debugf("\t - %s %s(%s)", method.ReturnType, method.Name, method.Args)
	}

	return method, err
}

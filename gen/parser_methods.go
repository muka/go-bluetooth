package gen

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (g *ApiGroup) parseMethods(raw []byte) []Method {

	methods := make([]Method, 0)
	slices := make([][]byte, 0)

	re := regexp.MustCompile(`(?s)Methods(.+?)(Properties|Signals)`)
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
		method := g.parseMethod(methodRaw)
		methods = append(methods, method)
	}

	return methods
}

func (g *ApiGroup) parseMethod(raw []byte) Method {

	method := Method{}
	// log.Debugf("%s", raw)

	re := regexp.MustCompile(`[ \t]*(.*?) ?(\w+)\(([^)]*)\) ?(.*?)\n((?s).+)`)
	matches1 := re.FindAllSubmatch(raw, -1)

	// log.Debugf("matches1 %s", matches1)
	for _, matches2 := range matches1 {

		rtype := string(matches2[1])
		if len(rtype) > 7 && rtype[:7] == "Methods" {
			rtype = rtype[7:]
		}
		rtype = strings.Trim(rtype, " \t")
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
				arg := Arg{
					Type: strings.Trim(argsparts[0], " \t\n"),
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

	if method.Name != "" && g.debug {
		log.Debugf("\t - %s %s(%s)", method.ReturnType, method.Name, method.Args)
	}

	return method
}

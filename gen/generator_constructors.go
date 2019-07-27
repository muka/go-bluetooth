package gen

import (
	"fmt"
	"regexp"
	"strings"
)

var defaultService = "org.bluez"

func isDefaultService(s string) bool {
	return len(s) >= len(defaultService) && s[:len(defaultService)] == defaultService
}

func createConstructors(api Api) []Constructor {

	constructors := []Constructor{}

	// log.Debugf("-------------------------------------- %s", api.Interface)
	// log.Debugf("ObjectPath %s", api.ObjectPath)
	// log.Debugf("Interface %s", api.Interface)

	apiService := api.Service
	if apiService != "" {
		apiService = strings.Split(apiService, " ")[0]
	}

	if !isDefaultService(apiService) {

		// log.Debugf("Service %s", api.Service)

		// case 1
		// unique name (Target role)
		// org.bluez (Controller role)
		if strings.Contains(apiService, "\n") {

			re := regexp.MustCompile(`(.+) \((.+?) .+\)`)
			matches1 := re.FindAllSubmatch([]byte(apiService), -1)

			for _, m1 := range matches1 {

				doc := ""
				srvc := strings.Trim(string(m1[1]), " \t")
				if !isDefaultService(srvc) {
					doc = srvc
					srvc = ""
				}

				c := Constructor{
					Service: srvc,
					Role:    string(m1[2]),
					Docs: []string{
						"servicePath: " + doc,
					},
				}

				constructors = append(constructors, c)
			}
		} else {
			c := Constructor{
				Service: "",
				Role:    "",
				Docs: []string{
					"servicePath: " + api.Service,
				},
			}
			constructors = append(constructors, c)
		}
		// log.Debugf("%++v", constructors)
	} else {
		c := Constructor{
			Service:    apiService,
			Role:       "",
			ObjectPath: "",
			Args:       "",
		}
		constructors = append(constructors, c)
	}

	for i, c := range constructors {

		inspectObjectPath(api.ObjectPath, &c)

		args := []string{}
		if c.Service == "" {
			args = append(args, "servicePath string")
			c.Service = "servicePath"
		} else {
			c.Service = fmt.Sprintf(`"%s"`, c.Service)
		}

		if c.ObjectPath == "" {
			args = append(args, "objectPath string")
			c.ObjectPath = "objectPath"
		} else {
			c.ObjectPath = fmt.Sprintf(`"%s"`, c.ObjectPath)
		}

		c.Args = strings.Join(args, ", ")

		docs := []string{}
		for _, doc := range c.Docs {
			for _, d1 := range strings.Split(doc, "\n") {
				docs = append(docs, "// "+d1)
			}
		}
		c.ArgsDocs = strings.Join(docs, "\n")

		constructors[i] = c
	}

	// log.Debugf("constructors %++v", constructors)

	return constructors
}

func inspectObjectPath(objectPath string, c *Constructor) {

	if strings.Contains(objectPath, "\n") {

		// log.Debugf("ObjectPath: \n----\n%s\n\n-----", objectPath)

		anchor1 := "(Target role)"
		idx := strings.Index(objectPath, anchor1)
		if idx > -1 {

			target := objectPath[:idx]

			anchor2 := "(Controller role)"
			idx2 := strings.Index(objectPath, anchor2)
			controller := objectPath[idx+len(anchor1) : idx2]

			target = strings.Replace(strings.Trim(target, " \t\n"), "\n", "", -1)
			controller = strings.Replace(strings.Trim(controller, " \t\n"), "\n\t", "", -1)

			// log.Debugf("target %s", target)
			// log.Debugf("controller %s", controller)
			// log.Debugf("ROLE %s", c.Role)

			c.Docs = append(c.Docs, "objectPath: "+objectPath)

			if c.Role == "Target" {
				c.ObjectPath = target
			}

			if c.Role == "Controller" {
				c.ObjectPath = controller
			}

		}

		return
	}

	// freely definable
	if strings.Contains(objectPath, "freely") {
		c.ObjectPath = ""
		c.Docs = append(c.Docs, "objectPath: "+objectPath)
		return
	}

	// <application dependent>
	if strings.HasPrefix(objectPath, "<application") {
		c.Docs = append(c.Docs, "objectPath: "+objectPath)
		c.ObjectPath = ""
		return
	}

}

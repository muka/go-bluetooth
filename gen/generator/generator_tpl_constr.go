package generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/muka/go-bluetooth/gen/override"
	"github.com/muka/go-bluetooth/gen/types"
)

var defaultService = "org.bluez"

func isDefaultService(s string) bool {
	return len(s) >= len(defaultService) && s[:len(defaultService)] == defaultService
}

func createConstructors(api *types.Api) []types.Constructor {

	// log.Debugf("-------------------------------------- %s", api.Interface)

	constructors := []types.Constructor{}
	constructors = inspectServiceName(api.Service, constructors)
	constructors = inspectObjectPath(api.ObjectPath, constructors)

	for i, c := range constructors {

		args := []string{}
		if c.Service == "" {
			args = append(args, "servicePath string")
			c.Service = "servicePath"
		} else {
			c.Service = fmt.Sprintf(`"%s"`, c.Service)
		}

		if c.ObjectPath == "" {
			args = append(args, "objectPath dbus.ObjectPath")
			c.ObjectPath = "objectPath"
		} else {
			c.ObjectPath = fmt.Sprintf(`"%s"`, c.ObjectPath)
		}

		c.Args = strings.Join(args, ", ")

		docs := []string{}
		for _, doc := range c.Docs {
			for _, d1 := range strings.Split(doc, "\n") {
				docs = append(docs, "// - "+d1)
			}
		}
		c.ArgsDocs = "//\n// Args:\n" + strings.Join(docs, "\n")

		// target as the deafult initializer
		if c.Role == "Target" {
			c.Role = ""
		}

		constructors[i] = c
	}

	if overrides, hasOverride := override.GetConstructorsOverrides(api.Interface); hasOverride {

		for _, coverride := range overrides {

			// add new constructors which take an adapter as arg
			if coverride.AdapterAsArgument {

				for _, c1 := range constructors {
					// log.Debugf("------ oveerride %+v", c1)

					c := types.Constructor{
						Args:       "adapterID string",
						ArgsDocs:   "// adapterID: ID of an adapter eg. hci0",
						Docs:       c1.Docs,
						ObjectPath: `fmt.Sprintf("/org/bluez/%s", adapterID)`,
						Service:    c1.Service,
						Role:       "FromAdapterID",
					}
					constructors = append(constructors, c)
				}

			}

		}

	}

	// log.Debugf("constructors %++v", constructors)

	return constructors
}

func inspectServiceName(serviceName string, constructors []types.Constructor) []types.Constructor {

	// log.Debugf("ObjectPath %s", api.ObjectPath)
	// log.Debugf("Interface %s", api.Interface)

	apiService := serviceName
	if apiService != "" {
		apiService = strings.Split(apiService, " ")[0]
	}

	if !isDefaultService(apiService) {

		// log.Debugf("Service %s", apiService)

		// case 1
		// unique name (Target role)
		// org.bluez (Controller role)
		if strings.Contains(serviceName, "\n") {

			re := regexp.MustCompile(`(.+) \((.+?) .+\)`)
			matches1 := re.FindAllSubmatch([]byte(serviceName), -1)

			// log.Debugf("%s ----> %s", serviceName, matches1)

			for _, m1 := range matches1 {

				doc := ""
				srvc := strings.Trim(string(m1[1]), " \t")

				if !isDefaultService(srvc) {
					doc = srvc
					srvc = ""
				}

				docslist := []string{}
				if doc != "" {
					docslist = append(docslist, "servicePath: "+doc)
				}

				c := types.Constructor{
					Service: srvc,
					Role:    string(m1[2]),
					Docs:    docslist,
				}

				constructors = append(constructors, c)
			}
		} else {

			c := types.Constructor{
				Service: "",
				Role:    "",
				Docs: []string{
					"servicePath: " + serviceName,
				},
			}
			constructors = append(constructors, c)
		}
	} else {
		c := types.Constructor{
			Service:    apiService,
			Role:       "",
			ObjectPath: "",
			Args:       "",
			Docs:       []string{},
		}
		constructors = append(constructors, c)
	}

	return constructors
}

func inspectObjectPath(objectPath string, constructors []types.Constructor) []types.Constructor {

	constructors2 := []types.Constructor{}

	// log.Debugf("%d %s", len(constructors), objectPath)
	// log.Debugf("%+v", constructors)

	for _, c := range constructors {

		if strings.Contains(objectPath, "\n") {

			// log.Debugf("ObjectPath: \n----\n%s\n\n-----", objectPath)

			anchor1 := " (Target role)"
			idx := strings.Index(objectPath, anchor1)
			if idx > -1 {

				target := objectPath[:idx]

				anchor2 := "(Controller role)"
				idx2 := strings.Index(objectPath, anchor2)
				controller := objectPath[idx+len(anchor1) : idx2]

				target = strings.Replace(strings.Trim(target, " \t\n"), "\n", "", -1)
				controller = strings.Replace(strings.Trim(controller, " \t\n"), "\n\t", "", -1)

				// if Role is set apply a objectPath
				if c.Role == "Target" {
					c.ObjectPath = ""
					c.Docs = append(c.Docs, "objectPath: "+target)
				}

				if c.Role == "Controller" {
					c.ObjectPath = ""
					c.Docs = append(c.Docs, "objectPath: "+controller)
				}

				// if no Role, create a contructor for each objectPath
				if c.Role == "" {

					if controller != "" {
						c1 := c
						c1.Role = "Controller"
						c1.ObjectPath = ""
						c1.Docs = append(c1.Docs, "objectPath: "+controller)
						constructors2 = append(constructors2, c1)
					}

					if target != "" {
						c1 := c
						c1.Role = "Target"
						c1.ObjectPath = ""
						c1.Docs = append(c1.Docs, "objectPath: "+target)
						constructors2 = append(constructors2, c1)
					}

					continue
				}

			}

			constructors2 = append(constructors2, c)
			continue
		}

		defaultObjectPath := "/org/bluez"

		isDefaultPath := len(objectPath) >= len(defaultObjectPath) && objectPath[:len(defaultObjectPath)] == defaultObjectPath
		pathHasVariables := strings.Contains(objectPath, "{")

		// log.Debugf("%s %t %t", objectPath, isDefaultPath, pathHasVariables)

		if !isDefaultPath || pathHasVariables {
			c.ObjectPath = ""
			c.Docs = append(c.Docs, "objectPath: "+objectPath)
		} else {
			c.ObjectPath = objectPath
		}

		// log.Debugf("----> %++v", c)
		constructors2 = append(constructors2, c)
	}

	return constructors2
}

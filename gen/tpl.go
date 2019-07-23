package gen

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

type BluezError struct {
	Name  string
	Error string
}

type BluezErrors struct {
	List []BluezError
}

func RootTemplate(filename string, api ApiGroup) error {

	fw, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %s", err)
	}

	tmpl := template.Must(template.ParseFiles("gen/tpl/root.go.tpl"))
	err = tmpl.Execute(fw, api)
	if err != nil {
		return fmt.Errorf("write tpl: %s", err)
	}
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

	tmpl := template.Must(template.ParseFiles("gen/tpl/errors.go.tpl"))

	err = tmpl.Execute(fw, errorsList)
	if err != nil {
		return fmt.Errorf("tpl write: %s", err)
	}

	return nil
}

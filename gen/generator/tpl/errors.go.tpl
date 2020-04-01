// Code generated DO NOT EDIT

package profile

import (
	"github.com/godbus/dbus"
)

var (
{{ range .List }}
	// {{.Name}} map to org.bluez.Error.{{.Name}}
	Err{{.Name}} = dbus.Error{
		Name: "org.bluez.Error.{{.Name}}",
		Body: []interface{}{"{{.Name}}"},
	}
{{ end }}
)

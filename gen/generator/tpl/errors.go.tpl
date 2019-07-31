// WARNING: generated code, do not edit
package profile

import (
	"github.com/godbus/dbus"
)

{{ range .List }}
// {{.Name}} map to org.bluez.Error.{{.Name}}
var {{.Name}} = dbus.NewError("org.bluez.Error.{{.Name}}", nil)
{{ end }}

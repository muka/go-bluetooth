package profile


const (
	OrgBluezInterface = "org.bluez"
{{ range .Interfaces }}
	//{{.InterfaceName}}Interface interface for {{.InterfaceName}}
	{{.InterfaceName}}Interface = "{{.Interface}}"
{{ end }}
)

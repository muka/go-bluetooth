package {{.Package}}
{{$InterfaceName := .InterfaceName}}

// {{.InterfaceName}} {{.Api.Title}}
{{.Api.Description}}
type {{.InterfaceName}} struct {
	client     *bluez.Client
	Properties *{{.InterfaceName}}Properties
}

// {{.InterfaceName}}Properties contains the exposed properties of an interface
type {{.InterfaceName}}Properties struct {
{{ range .Properties }}
	// {{.Property.Name}} {{.Property.Docs}}
	{{.Property.Name}} {{.Property.Type}}
{{end}}
}

// Close the connection
func (a *{{.InterfaceName}}) Close() {
	a.client.Disconnect()
}

//GetProperties load all available properties
func (a *{{.InterfaceName}}) GetProperties() (*{{.InterfaceName}}Properties, error) {
	err := a.client.GetProperties(a.Properties)
	return a.Properties, err
}

//SetProperty set a property
func (a *{{.InterfaceName}}) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

{{range .Methods}}
//{{.Name}} {{.Docs}}
func (a *{{$InterfaceName}}) {{.Name}}({{.ArgsList}}) {{.Method.ReturnType}} {
	return a.client.Call("{{.Name}}", 0).Store()
}
{{end}}

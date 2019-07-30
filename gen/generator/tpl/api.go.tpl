//
// WARNING: generated code, do not edit!
//
package {{.Package}}
{{$InterfaceName := .InterfaceName}}
{{$ExposeProperties := .ExposeProperties}}

{{.Imports}}

var {{.InterfaceName}}Interface = "{{.Api.Interface}}"

{{range .Constructors}}
// New{{$InterfaceName}}{{.Role}} create a new instance of {{$InterfaceName}}
{{.ArgsDocs}}
func New{{$InterfaceName}}{{.Role}}({{.Args}}) (*{{$InterfaceName}}, error) {
	a := new({{$InterfaceName}})
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  {{.Service}},
			Iface: {{$InterfaceName}}Interface,
			Path:  {{.ObjectPath}},
			Bus:   bluez.SystemBus,
		},
	)
	{{if $ExposeProperties }}
	a.Properties = new({{$InterfaceName}}Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	{{end}}
	return a, nil
}
{{end}}

// {{.InterfaceName}} {{.Api.Title}}
{{.Api.Description}}
type {{.InterfaceName}} struct {
	client     *bluez.Client
	Properties *{{.InterfaceName}}Properties
}

// {{.InterfaceName}}Properties contains the exposed properties of an interface
type {{.InterfaceName}}Properties struct {
	lock sync.RWMutex `dbus:"ignore"`
{{ range .Properties }}
	// {{.Property.Name}} {{.Property.Docs}}
	{{.Property.Name}} {{.Property.Type}}
{{end}}
}

func (p *{{.InterfaceName}}Properties) Lock() {
	p.lock.Lock()
}

func (p *{{.InterfaceName}}Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *{{.InterfaceName}}) Close() {
	a.client.Disconnect()
}

{{if .ExposeProperties }}
// ToMap convert a {{.InterfaceName}}Properties to map
func (a *{{.InterfaceName}}Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// GetProperties load all available properties
func (a *{{.InterfaceName}}) GetProperties() (*{{.InterfaceName}}Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *{{.InterfaceName}}) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *{{.InterfaceName}}) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *{{.InterfaceName}}) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *{{.InterfaceName}}) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}
{{end}}

{{range .Methods}}
//{{.Name}} {{.Docs}}
func (a *{{$InterfaceName}}) {{.Name}}({{.ArgsList}}) {{.Method.ReturnType}} {
	{{if .SingleReturn}}
	return a.client.Call("{{.Name}}", 0, {{.ParamsList}}).Store()
	{{else}}
	{{.ReturnVarsDefinition}}
	err := a.client.Call("{{.Name}}", 0, {{.ParamsList}}).Store({{.ReturnVarsRefs}})
	return {{.ReturnVarsList}}, err	{{end}}
}
{{end}}

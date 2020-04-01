// Code generated DO NOT EDIT

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
			Path:  dbus.ObjectPath({{.ObjectPath}}),
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

/*
{{.InterfaceName}} {{.Api.Title}}
{{.Api.Description}}
*/
type {{.InterfaceName}} struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*{{.InterfaceName}}Properties
	watchPropertiesChannel chan *dbus.Signal
}

// {{.InterfaceName}}Properties contains the exposed properties of an interface
type {{.InterfaceName}}Properties struct {
	lock sync.RWMutex `dbus:"ignore"`
{{ range .Properties }}
	/*
	{{.Property.Name}} {{.Property.Docs}}
	*/
	{{.Property.Name}} {{.Property.Type}}
{{end}}
}

//Lock access to properties
func (p *{{.InterfaceName}}Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *{{.InterfaceName}}Properties) Unlock() {
	p.lock.Unlock()
}

{{ range .Properties }}

{{ if not .ReadOnly }}
// Set{{.Property.Name}} set {{.Property.Name}} value
func (a *{{$InterfaceName}}) Set{{.Property.Name}}(v {{.RawType}}) error {
	return a.SetProperty("{{.Property.Name}}", v)
}
{{end}}

{{if not .WriteOnly }}
// Get{{.Property.Name}} get {{.Property.Name}} value
func (a *{{$InterfaceName}}) Get{{.Property.Name}}() ({{.RawType}}, error) {
	v, err := a.GetProperty("{{.Property.Name}}")
	if err != nil {
		return {{.RawTypeInitializer}}, err
	}
	return v.Value().({{.RawType}}), nil
}
{{end}}
{{end}}

// Close the connection
func (a *{{.InterfaceName}}) Close() {
	{{if $ExposeProperties }}
	a.unregisterPropertiesSignal()
	{{end}}
	a.client.Disconnect()
}

// Path return {{.InterfaceName}} object path
func (a *{{.InterfaceName}}) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return {{.InterfaceName}} dbus client
func (a *{{.InterfaceName}}) Client() *bluez.Client {
	return a.client
}

// Interface return {{.InterfaceName}} interface
func (a *{{.InterfaceName}}) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *{{.InterfaceName}}) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
}

{{if .ExposeProperties }}
// ToMap convert a {{.InterfaceName}}Properties to map
func (a *{{.InterfaceName}}Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an {{.InterfaceName}}Properties
func (a *{{.InterfaceName}}Properties) FromMap(props map[string]interface{}) (*{{.InterfaceName}}Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an {{.InterfaceName}}Properties
func (a *{{.InterfaceName}}Properties) FromDBusMap(props map[string]dbus.Variant) (*{{.InterfaceName}}Properties, error) {
	s := new({{.InterfaceName}}Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *{{.InterfaceName}}) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *{{.InterfaceName}}) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *{{.InterfaceName}}) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
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

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *{{.InterfaceName}}) GetPropertiesSignal() (chan *dbus.Signal, error) {

	if a.propertiesSignal == nil {
		s, err := a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		a.propertiesSignal = s
	}

	return a.propertiesSignal, nil
}

// Unregister for changes signalling
func (a *{{.InterfaceName}}) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *{{.InterfaceName}}) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *{{.InterfaceName}}) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}

{{end}}

{{range .Methods}}
/*
{{.Name}} {{.Docs}}
*/
func (a *{{$InterfaceName}}) {{.Name}}({{.ArgsList}}) {{.Method.ReturnType}} {
	{{if .SingleReturn}}
	return a.client.Call("{{.Name}}", 0, {{.ParamsList}}).Store()
	{{else}}
	{{.ReturnVarsDefinition}}
	err := a.client.Call("{{.Name}}", 0, {{.ParamsList}}).Store({{.ReturnVarsRefs}})
	return {{.ReturnVarsList}}, err	{{end}}
}
{{end}}

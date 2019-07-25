package gen

type Flag int

const (
	FlagReadOnly     Flag = 1
	FlagWriteOnly    Flag = iota
	FlagReadWrite    Flag = iota
	FlagExperimental Flag = iota
)

type Arg struct {
	Type string
	Name string
}

type Method struct {
	Name       string
	ReturnType string
	Args       []Arg
	Errors     []string
	Docs       string
}

type Property struct {
	Name  string
	Type  string
	Docs  string
	Flags []Flag
}

type ApiGroup struct {
	FileName    string
	Name        string
	Description string
	Api         []Api
	debug       bool
}

type Api struct {
	Title       string
	Description string
	Service     string
	Interface   string
	ObjectPath  string
	Methods     []Method
	// those are currently avail only in health-api
	Signals    []Method
	Properties []Property
}

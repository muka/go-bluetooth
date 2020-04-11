package types

type ApiGroup struct {
	FileName    string
	Name        string
	Description string
	Api         []*Api
	debug       bool
}

type Api struct {
	Title       string
	Description string
	Service     string
	Interface   string
	ObjectPath  string
	Methods     []*Method
	// those are currently avail only in health-api
	Signals    []*Method
	Properties []*Property
}

type Flag int

const (
	FlagReadOnly     Flag = 1
	FlagWriteOnly    Flag = 2
	FlagReadWrite    Flag = 3
	FlagExperimental Flag = 4
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

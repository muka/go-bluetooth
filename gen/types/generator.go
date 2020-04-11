package types

type BluezError struct {
	Name  string
	Error string
}

type BluezErrors struct {
	List []BluezError
}

type MethodDoc struct {
	*Method
	ArgsList             string
	ParamsList           string
	SingleReturn         bool
	ReturnVarsDefinition string
	ReturnVarsRefs       string
	ReturnVarsList       string
}

type InterfaceDoc struct {
	Title     string
	Name      string
	Interface string
}

type InterfacesDoc struct {
	Interfaces []InterfaceDoc
}

type PropertyDoc struct {
	*Property
	RawType            string
	RawTypeInitializer string
	ReadOnly           bool
	WriteOnly          bool
	ReadWrite          bool
}

type ApiGroupDoc struct {
	*ApiGroup
	Package string
}

type ApiDoc struct {
	Api              *Api
	InterfaceName    string
	Package          string
	Properties       []PropertyDoc
	Methods          []MethodDoc
	Imports          string
	Constructors     []Constructor
	ExposeProperties bool
}

type Constructor struct {
	Service    string
	Role       string
	ObjectPath string
	Args       string
	ArgsDocs   string
	Docs       []string
}

package generator

import (
	"embed"
	"strings"
	"text/template"

	"github.com/muka/go-bluetooth/gen/types"
)

//go:embed tpl/*.tpl
var _tplFS embed.FS

var tpl = template.Must(
	template.New("").
		ParseFS(_tplFS, "tpl/*.tpl"),
)

// rename variable name to avoid collision with Go languages
func renameReserved(varname string) string {
	switch varname {
	case "type":
		return "type1"
	default:
		return varname
	}
}

func loadtpl(name string) *template.Template {
	res := tpl.Lookup(name + ".go.tpl")
	if res == nil {
		panic("not found template with name " + name)
	}
	return res
}

func prepareDocs(src string, skipFirstComment bool, leftpad int) string {
	return src
	// lines := strings.Split(src, "\n")
	// result := []string{}
	// // comment := "// "
	// comment := ""
	// prefixLen := leftpad + len(comment)
	// fmtt := fmt.Sprintf("%%%ds%%s", prefixLen)
	//
	// for _, line := range lines {
	// 	line = strings.Trim(line, " \t\r")
	// 	if len(line) == 0 {
	// 		continue
	// 	}
	//
	// 	result = append(result, fmt.Sprintf(fmtt, comment, line))
	// }
	// if skipFirstComment && len(result) > 0 && len(result[0]) > 3 {
	// 	result[0] = result[0][prefixLen:]
	// }
	// return strings.Join(result, "\n")
}

func getApiPackage(apiGroup *types.ApiGroup) string {
	apiName := strings.Replace(apiGroup.FileName, "-api.txt", "", -1)
	apiName = strings.Replace(apiName, "-", "_", -1)
	apiName = strings.Replace(apiName, " [experimental]", "", -1)
	return apiName
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

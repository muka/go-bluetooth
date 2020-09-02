package filters

import (
	"fmt"
	"os"
	"strings"
)

type FilterContext int

const (
	FilterFile   FilterContext = 0
	FilterApi    FilterContext = iota
	FilterMethod FilterContext = iota
)

const (
	ParamFileFilter   = "file_filter"
	ParamApiFilter    = "api_filter"
	ParamMethodFilter = "method_filter"
)

type Filter struct {
	Context FilterContext
	Value   string
}

func NewFilter(value string, context FilterContext) Filter {
	return Filter{context, value}
}

func extractFilters(param string, paramType FilterContext) []Filter {
	list := []Filter{}

	// parse from env vars
	rawFilters := strings.Split(os.Getenv(strings.ToUpper(param)), ",")
	for _, filter := range rawFilters {
		filter = strings.Trim(filter, " ")
		if len(filter) == 0 {
			continue
		}
		list = append(list, NewFilter(filter, paramType))
	}

	// parse from args
	if len(os.Args) > 1 {
		args := os.Args[1:]
		for _, arg := range args {
			if strings.Contains(arg, fmt.Sprintf("%s=", param)) {
				filters2 := strings.Split(strings.Split(arg, "=")[1], ",")
				for _, filter := range filters2 {
					filter = strings.Trim(filter, " ")
					if len(filter) == 0 {
						continue
					}

					list = append(list, NewFilter(filter, paramType))
				}
			}
		}
	}

	return list
}

func ParseCliFilters() []Filter {
	filters := []Filter{}
	filters = append(filters, extractFilters(ParamFileFilter, FilterFile)...)
	filters = append(filters, extractFilters(ParamApiFilter, FilterApi)...)
	filters = append(filters, extractFilters(ParamMethodFilter, FilterMethod)...)
	return filters
}

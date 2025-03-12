package stringutil

import (
	"strings"
	"text/template"
)

func LastIndexOf(slice []string, value string) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

func ExecuteTemplate(tmpl string, params any) string {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	var buf strings.Builder
	if err = t.Execute(&buf, params); err != nil {
		panic(err)
	}
	return buf.String()
}

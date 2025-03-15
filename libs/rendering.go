package libs

import (
	"io"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func RenderTemplate(tmpl string, data any) ([]byte, error) {

	WriteOutputf("Rendering template: %s\n", tmpl)

	textTpl := template.New("template").Funcs(sprig.FuncMap())

	tpl, err := textTpl.Parse(string(tmpl))
	if err != nil {
		return nil, err
	}

	var outBuffer strings.Builder
	outIO := io.MultiWriter(&outBuffer)

	err = tpl.Execute(outIO, data)
	if err != nil {
		return nil, err
	}

	return []byte(outBuffer.String()), nil
}

func MergeMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

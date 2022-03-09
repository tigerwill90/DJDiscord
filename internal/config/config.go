package config

import (
	"html/template"
	"io"
)

var t = template.Must(template.New("config.txt").Parse(tmpl))

func Generate(w io.Writer, opts *TemplateOption) error {
	return t.Execute(w, opts)
}

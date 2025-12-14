package templates

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var tmplFS embed.FS

func New(name string) *template.Template {
	funcMap := template.FuncMap{
		// you can add functions to use in templates here
	}
	return template.Must(
		template.New(name).Funcs(funcMap).ParseFS(
			tmplFS, "templates/"+name))
}

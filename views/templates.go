package views

import "html/template"

var tpl *template.Template

func Tpl() *template.Template {
	if tpl == nil {
		tpl = template.Must(template.ParseGlob("views/templates/*"))
	}

	return tpl
}

package views

import "html/template"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func Tpl() *template.Template {
	return tpl
}

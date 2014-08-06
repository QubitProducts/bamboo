package configuration

import (
//	"text/template"
)

type HAProxy struct {
	TemplatePath string
	OutputPath string
	ReloadCommand string
}

func (h HAProxy) WriteTemplate(apps interface {}, domainMapping interface {}) {
//	tpl := template.New(h.TemplatePath)
//	var renderedText *string

//	tpl.Execute(renderedText, apps, domainMapping)
}

func Reload() {
	// TODO call reload system call
}

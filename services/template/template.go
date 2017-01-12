package template

import (
	"bytes"
	"github.com/cloverstd/bamboo/services/service"
	"strings"
	"text/template"
)

func divide(a int, b int) int {
	return a / b
}

func multip(a int, b int) int {
	return a * b
}

func minus(input int, change int) int {
	return input - change
}

func plus(input int, change int) int {
	return input + change
}

func hasKey(data map[string]service.Service, appId string) bool {
	_, exists := data[appId]
	return exists
}

func getService(data map[string]service.Service, appId string) service.Service {
	serviceModel, _ := data[appId]
	return serviceModel
}

/*
	Returns string content of a rendered template
*/
func RenderTemplate(templateName string, templateContent string, data interface{}) (string, error) {
	funcMap := template.FuncMap{
		"divide":     divide,
		"multip":     multip,
		"minus":      minus,
		"plus":       plus,
		"hasKey":     hasKey,
		"getService": getService,
		"Split":      strings.Split,
		"Contains":   strings.Contains,
		"Join":       strings.Join,
		"Replace":    strings.Replace,
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower}

	tpl := template.Must(template.New(templateName).Funcs(funcMap).Parse(templateContent))

	strBuffer := new(bytes.Buffer)

	err := tpl.Execute(strBuffer, data)
	if err != nil {
		return "", err
	}

	return strBuffer.String(), nil
}

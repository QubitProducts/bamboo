package template

import (
	"bytes"
	"github.com/QubitProducts/bamboo/services/service"
	"github.com/QubitProducts/bamboo/services/marathon"
	"strings"
	"text/template"
)

func hasKey(data map[string]service.Service, appId string) bool {
	_, exists := data[appId]
	return exists
}

func getService(data map[string]service.Service, appId string) service.Service {
	serviceModel, _ := data[appId]
	return serviceModel
}

func getAppEnvValues(envVar string, appList marathon.AppList) (envValues []string) {
	for _, app := range appList {
		value, ok := app.Env[envVar]
		if ok {
			envValues = append(envValues, value)
		}
	}
	return unique(envValues)
}

func unique(strings []string) (uniqueStrings []string) {
	seen := make(map[string]bool)
	for _, s := range strings {
		if !seen[s] {
			uniqueStrings = append(uniqueStrings, s)	    
			seen[s] = true
		}
	}
	return uniqueStrings
}

/*
	Returns string content of a rendered template
*/
func RenderTemplate(templateName string, templateContent string, data interface{}) (string, error) {
	funcMap := template.FuncMap{
		"hasKey":          hasKey,
		"getService":      getService,
		"getAppEnvValues": getAppEnvValues,
		"Split":      	   strings.Split,
		"Contains":   	   strings.Contains,
		"Join":       	   strings.Join,
		"Replace":    	   strings.Replace,
		"ToUpper":    	   strings.ToUpper,
		"ToLower":    	   strings.ToLower,
		"unique":     	   unique}

	tpl := template.Must(template.New(templateName).Funcs(funcMap).Parse(templateContent))

	strBuffer := new(bytes.Buffer)

	err := tpl.Execute(strBuffer, data)
	if err != nil {
		return "", err
	}

	return strBuffer.String(), nil
}

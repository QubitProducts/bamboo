package writer

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

/*
	Writes template a given output file path
*/
func WriteTemplate(templatePath string, outputFilePath string, data interface{}) error {

	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return err
	}

	content, errRender := RenderTemplate(templatePath, string(templateContent), data)
	if errRender != nil {
		return errRender
	}

	return ioutil.WriteFile(outputFilePath, []byte(content), 0666)
}


func hasKey(data map[string]string, appId string) bool {
	_, exists := data[appId]
	return exists
}

func getValue(data map[string]string, appId string) string {
	value, _ := data[appId]
	return value
}

/*
	Returns string content of a rendered template
*/
func RenderTemplate(templateName string, templateContent string, data interface{}) (string, error) {
	funcMap := template.FuncMap{ "hasKey": hasKey, "getValue": getValue }

	tpl := template.Must(template.New(templateName).Funcs(funcMap).Parse(templateContent))

	strBuffer := new(bytes.Buffer)

	err := tpl.Execute(strBuffer, data)
	if err != nil {
		return "", err
	}

	return strBuffer.String(), nil
}


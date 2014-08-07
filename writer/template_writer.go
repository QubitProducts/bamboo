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

/*
	Returns string content of a rendered template
*/
func RenderTemplate(templateName string, templateContent string, data interface{}) (string, error) {
	tpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return "", err
	}

	strBuffer := new(bytes.Buffer)

	err = tpl.Execute(strBuffer, data)
	if err != nil {
		return "", err
	}

	return strBuffer.String(), nil
}

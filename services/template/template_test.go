package template

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/QubitProducts/bamboo/services/marathon"
)

func TestTemplateWriter(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		templateContent := "{{.id}} {{.domain}}"
		params := map[string]string{"id": "app", "domain": "example.com"}

		Convey("should render template as string", func() {
			content, _ := RenderTemplate(templateName, templateContent, params)
			So(content, ShouldEqual, "app example.com")
		})
	})
}

type TemplateData struct {
        Apps     marathon.AppList
}

func TestTemplateGetAppEnvValuesFunction(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		apps := marathon.AppList{}
		apps = append(apps, marathon.App{
		     Id: "app1",
		     Env: map[string]string{
		     	  "BAMBOO_TCP_PORT": "10001",
		     },
		     })
		apps = append(apps, marathon.App{
		     Id: "app2",
		     Env: map[string]string{
		     	  "BAMBOO_TCP_PORT": "10001",
		     },
		     })
		apps = append(apps, marathon.App{
		     Id: "app3",
		     Env: map[string]string{
		     	  "BAMBOO_TCP_PORT": "10002",
		     },
		     })
		apps = append(apps, marathon.App{
		     Id: "app4",
		     Env: map[string]string{
		     	  "BAMBOO_TCP_PORT": "10001",
		     },
		     })
		apps = append(apps, marathon.App{
		     Id: "app5",
		     Env: map[string]string{
		     	  "BAMBOO_TCP_PORT": "10003",
		     },
		     })
		apps = append(apps, marathon.App{
		     Id: "app6",
		     Env: map[string]string{
		     	  "BAMBOO_TCP_PORT": "10002",
		     },
		     })

		templateData := TemplateData{Apps: apps}
		templateContent := "{{ range $port := .Apps | getAppEnvValues \"BAMBOO_TCP_PORT\" }}{{ $port }} {{ end }}"
		Convey("should list the three unique ports", func() {
			content, _ := RenderTemplate(templateName, templateContent, templateData)
			So(content, ShouldEqual, "10001 10002 10003 ")
		})
	})
}

func TestTemplateSplitFunction(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		params := map[string]string{"id": "app", "domain": "example.com"}

		templateContent := "{{.id}}{{range Split .domain \".\"}} {{.}}{{end}}"
		Convey("should split domain as two different strings", func() {
			content, _ := RenderTemplate(templateName, templateContent, params)
			So(content, ShouldEqual, "app example com")
		})
	})
}

func TestTemplateContainsFunction(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		params := map[string]string{"ip": "10.12.12.15", "subnet": "10.12"}

		templateContent := "{{if Contains .ip .subnet}}true{{end}}"
		Convey("should verify if contains ip in subnet", func() {
			content, _ := RenderTemplate(templateName, templateContent, params)
			So(content, ShouldEqual, "true")
		})
	})
}

func TestTemplateJoinFunction(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		domains := []string{"example.com", "example.net"}
		params := map[string][]string{"domains": domains}

		templateContent := "{{Join .domains \", \"}}"
		Convey("should create a list of two domains separated by comma", func() {
			content, _ := RenderTemplate(templateName, templateContent, params)
			So(content, ShouldEqual, "example.com, example.net")
		})
	})
}

func TestTemplateReplaceFunction(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		params := map[string]string{"domain": "example.com"}

		templateContent := "{{Replace .domain \"com\" \"net\" 1}}"
		Convey("should replace example.com into example.net", func() {
			content, _ := RenderTemplate(templateName, templateContent, params)
			So(content, ShouldEqual, "example.net")
		})
	})
}

func TestTemplateToUpperFunction(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		params := map[string]string{"id": "app", "domain": "example.com"}

		templateContent := "{{.id}} {{ToUpper .domain}}"
		Convey("should transform example.com upper case", func() {
			content, _ := RenderTemplate(templateName, templateContent, params)
			So(content, ShouldEqual, "app EXAMPLE.COM")
		})
	})
}

func TestTemplateToLowerFunction(t *testing.T) {
	Convey("#RenderTemplate", t, func() {
		templateName := "templateName"
		params := map[string]string{"domain": "EXAMPLE.COM"}

		templateContent := "{{ToLower .domain}}"
		Convey("should transform example.com lower case", func() {
			content, _ := RenderTemplate(templateName, templateContent, params)
			So(content, ShouldEqual, "example.com")
		})
	})
}

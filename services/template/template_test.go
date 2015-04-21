package template

import (
	. "github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"testing"
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
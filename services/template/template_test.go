package template

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
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

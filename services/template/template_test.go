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

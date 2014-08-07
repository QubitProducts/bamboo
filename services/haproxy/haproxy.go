package haproxy

import(
	"bamboo/services/marathon"

	conf "bamboo/configuration"
	writer "bamboo/writer"
)

func WriteHAProxyConfig(haproxyConf conf.HAProxy, apps []marathon.App, services map[string]string) error {

	// data for rendering template
	data := struct {
			Apps     []marathon.App
			Services map[string]string
		}{apps, services}

	return writer.WriteTemplate(haproxyConf.TemplatePath, haproxyConf.OutputPath, data)
}

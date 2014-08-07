package haproxy

import(
	"bamboo/services/marathon"

	conf "bamboo/configuration"
	writer "bamboo/writer"
	"bamboo/services/cmd"
)


type templateData struct {
	Apps []marathon.App
	Services map[string]string
}

func WriteHAProxyConfig(haproxyConf conf.HAProxy, data interface {}) error {
	return writer.WriteTemplate(haproxyConf.TemplatePath, haproxyConf.OutputPath, data)
}

func GetTemplateData() interface {} {
	config := cmd.GetConfiguration()

	apps, _ := marathon.FetchApps(config.Marathon.Endpoint)
	//	services, _ := domain.FetchAll(config.ServicesMapping.Zookeeper)
	services := map[string]string{}

	return templateData{ apps, services }
}

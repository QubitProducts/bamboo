package haproxy

import (
	"github.com/samuel/go-zookeeper/zk"

	conf "bamboo/configuration"
	"bamboo/services/marathon"
	"bamboo/services/domain"
	"bamboo/writer"
)

type templateData struct {
	Apps    []marathon.App
	Services map[string]string
}

func WriteHAProxyConfig(haproxyConf conf.HAProxy, data interface{}) error {
	return writer.WriteTemplate(haproxyConf.TemplatePath, haproxyConf.OutputPath, data)
}

func GetTemplateData(config conf.Configuration, conn *zk.Conn) interface{} {

	apps, _ := marathon.FetchApps(config.Marathon.Endpoint)
	services, _ := domain.All(conn, config.DomainMapping.Zookeeper)

	return templateData{ apps, services }
}

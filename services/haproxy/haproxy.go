package haproxy

import (
	"github.com/samuel/go-zookeeper/zk"

	conf "github.com/QubitProducts/bamboo/configuration"
	"github.com/QubitProducts/bamboo/services/marathon"
	"github.com/QubitProducts/bamboo/services/domain"
	"github.com/QubitProducts/bamboo/writer"
)

type templateData struct {
	Apps    []marathon.App
	Services map[string]string
}

func WriteHAProxyConfig(haproxyConf conf.HAProxy, data interface{}) error {
	return writer.WriteTemplate(haproxyConf.TemplatePath, haproxyConf.OutputPath, data)
}

func GetTemplateData(config *conf.Configuration, conn *zk.Conn) interface{} {

	apps, _ := marathon.FetchApps(config.Marathon.Endpoint)
	services, _ := domain.All(conn, config.DomainMapping.Zookeeper)

	return templateData{ apps, services }
}

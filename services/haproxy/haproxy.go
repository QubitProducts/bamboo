package haproxy

import (
	conf "github.com/cloverstd/bamboo/configuration"
	"github.com/cloverstd/bamboo/services/marathon"
	"github.com/cloverstd/bamboo/services/service"
)

type templateData struct {
	Apps     marathon.AppList
	Services map[string]service.Service
}

func GetTemplateData(config *conf.Configuration, storage service.Storage) (*templateData, error) {

	apps, err := marathon.FetchApps(config.Marathon, config)

	if err != nil {
		return nil, err
	}

	services, err := storage.All()
	if err != nil {
		return nil, err
	}

	byName := make(map[string]service.Service)
	for _, service := range services {
		byName[service.Id] = service
	}

	return &templateData{apps, byName}, nil
}

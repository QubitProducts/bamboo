package domain

import (
	conf "bamboo/configuration"
)



/*
	Get an  key value
 */
func FetchAll(zkConf conf.Zookeeper) (map[string]string, error) {

	return map[string]string{}, nil
}

func Create(zkConf conf.Zookeeper, appName string, domainValue string) error {
	return nil
}

func Update(zkConf conf.Zookeeper, appName string, domainValue string) error {
	return nil
}

func Delete(zkConf conf.Zookeeper, appName string) error {
	return nil
}

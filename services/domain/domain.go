package domain

import (
	"github.com/samuel/go-zookeeper/zk"
	conf "bamboo/configuration"
)

/*
	Get an  key value
 */
func All(conn *zk.Conn, zkConf conf.Zookeeper) (map[string]string, error) {
//	zkConf.Path

	return map[string]string{}, nil
}


func Create(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) error {
//	createRq := CreateRequest{ Path: zkConf.Path }
//	ops := zk.MultiOps{ []CreateRequest: createRq }
	return nil
}

func Put(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) error {
	return nil
}

func Delete(conn *zk.Conn, zkConf conf.Zookeeper, appId string) error {
	return nil
}
//
//func zkConfig() {
//	return cmd.GetConfiguration().DomainMapping.Zookeeper
//}



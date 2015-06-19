package service

import (
	"github.com/QubitProducts/bamboo/Godeps/_workspace/src/github.com/samuel/go-zookeeper/zk"
	conf "github.com/QubitProducts/bamboo/configuration"
	"net/url"
	"strings"
)

type Service struct {
	Id  string `param:"id"`
	Acl string `param:"acl"`
}

func All(conn *zk.Conn, zkConf conf.Zookeeper) (map[string]Service, error) {

	err := ensurePathExists(conn, zkConf.Path)
	if err != nil {
		return nil, err
	}

	services := map[string]Service{}
	keys, _, err2 := conn.Children(zkConf.Path)

	if err2 != nil {
		return nil, err2
	}

	for _, childPath := range keys {
		bite, _, e := conn.Get(zkConf.Path + "/" + childPath)
		if e != nil {
			return nil, e
			break
		}
		appId, _ := unescapeSlashes(childPath)
		services[appId] = Service{Id: appId, Acl: string(bite)}
	}
	return services, nil
}

/*
   Read ZK ACL:
   http://zookeeper.apache.org/doc/trunk/zookeeperProgrammers.html#sc_ACLPermissions
*/
func Create(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) (string, error) {
	path := concatPath(zkConf.Path, validateAppId(appId))
	resPath, err := conn.Create(path, []byte(domainValue), 0, defaultACL())
	if err != nil {
		return "", err
	}

	return resPath, nil
}

func Put(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) (*zk.Stat, error) {
	path := concatPath(zkConf.Path, validateAppId(appId))
	err := ensurePathExists(conn, path)
	if err != nil {
		return nil, err
	}

	stats, err := conn.Set(path, []byte(domainValue), -1)
	if err != nil {
		return nil, err
	}

	// Force triger an event on parent
	conn.Set(zkConf.Path, []byte{}, -1)

	return stats, nil
}

func Delete(conn *zk.Conn, zkConf conf.Zookeeper, appId string) error {
	path := concatPath(zkConf.Path, validateAppId(appId))
	return conn.Delete(path, -1)
}

func concatPath(parentPath string, appId string) string {
	return parentPath + "/" + escapeSlashes(appId)
}

func ensurePathExists(conn *zk.Conn, path string) error {
	pathExists, _, _ := conn.Exists(path)
	if pathExists {
		return nil
	}

	_, err := conn.Create(path, []byte{}, 0, defaultACL())
	if err != nil {
		return err
	}

	return nil
}

func defaultACL() []zk.ACL {
	return []zk.ACL{zk.ACL{Perms: zk.PermAll, Scheme: "world", ID: "anyone"}}
}

func validateAppId(appId string) string {
	if strings.HasPrefix(appId, "/") {
		return appId
	} else {
		return "/" + appId
	}
}

func escapeSlashes(id string) string {
	return url.QueryEscape(id)
}

func unescapeSlashes(id string) (string, error) {
	return url.QueryUnescape(id)
}

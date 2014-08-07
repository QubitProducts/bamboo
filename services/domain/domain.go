package domain

import (
	"fmt"

	"github.com/samuel/go-zookeeper/zk"
	conf "bamboo/configuration"
)


func All(conn *zk.Conn, zkConf conf.Zookeeper) (map[string]string, error) {

	fmt.Println(zkConf.Path)

	// TODO: ensure nested state is created
	// check path exists
	err := ensurePathExists(conn, zkConf.Path)
	if err != nil { return nil, err }

	domains := map[string]string{}
	keys, _, err2 := conn.Children(zkConf.Path)

	if err2 != nil { return nil, err2 }

	for _, childPath := range keys {
		bite, _, e := conn.Get(childPath)
		if e != nil {
			return nil, e
			break
		}
		domains[childPath] = string(bite)
	}
	return domains, nil
}


/*
   Read about ZK ACL:
   http://zookeeper.apache.org/doc/trunk/zookeeperProgrammers.html#sc_ACLPermissions
*/

func Create(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) (string, error) {
	path := concatPath(zkConf.Path, appId)
	resPath, err := conn.Create(path, []byte(domainValue), 0, defaultACL())
	if err != nil { return "", err }

	return resPath, nil
}

func Put(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) (*zk.Stat, error) {
	path := concatPath(zkConf.Path, appId)
	// use default version
	stats, err := conn.Set(path, []byte(domainValue), -1)

	if err != nil { return nil, err }

	return stats, nil
}

func Delete(conn *zk.Conn, zkConf conf.Zookeeper, appId string) error {
	path := concatPath(zkConf.Path, appId)
	return conn.Delete(path, -1)
}

func concatPath(parentPath string, appId string) string {
	return parentPath + "/" + appId
}

func ensurePathExists(conn *zk.Conn, path string) error {
	pathExists, _, _ := conn.Exists(path)
	if pathExists { return nil }

	_, err := conn.Create(path, []byte{}, 0, defaultACL())
	if err != nil { return err }

	return nil
}

func defaultACL() []zk.ACL {
	return []zk.ACL{ zk.ACL { Perms: zk.PermAll, Scheme: "world", ID: "anyone" }}
}

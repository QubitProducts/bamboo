package domain

import (
	conf "bamboo/configuration"
	"github.com/samuel/go-zookeeper/zk"
)

func All(conn *zk.Conn, zkConf conf.Zookeeper) (map[string]string, error) {

	// TODO: ensure nested state is created
	// check path exists
	err := ensurePathExists(conn, zkConf.Path)
	if err != nil {
		return nil, err
	}

	domains := map[string]string{}
	keys, _, err2 := conn.Children(zkConf.Path)

	if err2 != nil {
		return nil, err2
	}

	for _, childPath := range keys {
		bite, _, e := conn.Get(concatPath(zkConf.Path, childPath))
		if e != nil {
			return nil, e
			break
		}
		domains[childPath] = string(bite)
	}
	return domains, nil
}

/*
   Read ZK ACL:
   http://zookeeper.apache.org/doc/trunk/zookeeperProgrammers.html#sc_ACLPermissions
*/
func Create(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) (string, error) {
	path := concatPath(zkConf.Path, appId)
	resPath, err := conn.Create(path, []byte(domainValue), 0, defaultACL())
	if err != nil {
		return "", err
	}

	return resPath, nil
}


func Put(conn *zk.Conn, zkConf conf.Zookeeper, appId string, domainValue string) (*zk.Stat, error) {
	path := concatPath(zkConf.Path, appId)
	stats, err := conn.Set(path, []byte(domainValue), -1)

	if err != nil {
		return nil, err
	}
	// Force triger an event on parent
	conn.Set(zkConf.Path, []byte{}, -1)

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

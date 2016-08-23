package service

import (
	conf "github.com/QubitProducts/bamboo/configuration"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net/url"
)

type ZKStorage struct {
	conn *zk.Conn
	conf conf.Zookeeper
	acl  []zk.ACL
}

func NewZKStorage(conn *zk.Conn, conf conf.Zookeeper) (s *ZKStorage, err error) {
	s = &ZKStorage{
		conn: conn,
		conf: conf,
		acl:  defaultACL(),
	}
	err = s.ensurePathExists()
	return s, err
}

func (z *ZKStorage) All() (services []Service, err error) {
	err = z.ensurePathExists()
	if err != nil {
		return
	}

	keys, _, err := z.conn.Children(z.conf.Path)
	if err != nil {
		return
	}

	services = make([]Service, 0, len(keys))
	for _, childPath := range keys {
		body, _, err := z.conn.Get(z.conf.Path + "/" + childPath)
		if err != nil {
			return nil, err
		}

		path, err := unescapePath(childPath)
		if err != nil {
			return nil, err
		}

		// We tolerate being unable to decode a service body, as may be new version running simultaneously
		repr, err := ParseServiceRepr(body, path)
		if err != nil {
			log.Printf("Failed to parse service at %v: %v", path, err)
			continue
		}

		services = append(services, repr.Service())
	}

	return
}

func (z *ZKStorage) Upsert(service Service) (err error) {
	repr := MakeV2ServiceRepr(service)

	body, err := repr.Serialize()
	if err != nil {
		return
	}

	err = z.ensurePathExists()
	if err != nil {
		return err
	}

	path := z.servicePath(service.Id)

	ok, _, err := z.conn.Exists(path)
	if err != nil {
		return
	}

	if ok {
		_, err = z.conn.Set(path, body, -1)
		if err != nil {
			log.Print("Failed to set path", err)
			return
		}

		// Trigger an event on the parent
		_, err = z.conn.Set(z.conf.Path, []byte{}, -1)
		if err != nil {
			log.Print("Failed to trigger event on parent", err)
			err = nil
		}

	} else {
		_, err = z.conn.Create(path, body, 0, z.acl)
		if err != nil {
			log.Print("Failed to set create", err)
			return
		}
	}
	return
}

func (z *ZKStorage) Delete(serviceId string) error {
	path := z.servicePath(serviceId)
	return z.conn.Delete(path, -1)
}

func (z *ZKStorage) servicePath(id string) string {
	return z.conf.Path + "/" + escapePath(id)
}

func (z *ZKStorage) ensurePathExists() error {
	pathExists, _, _ := z.conn.Exists(z.conf.Path)
	if pathExists {
		return nil
	}

	// This is a fairly rare, and fairly critical, operation, so I'm going to be verbose
	log.Print("Creating base zk path", z.conf.Path)
	_, err := z.conn.Create(z.conf.Path, []byte{}, 0, z.acl)
	if err != nil {
		log.Print("Failed to create base zk path", err)
	}

	return err
}

func defaultACL() []zk.ACL {
	return []zk.ACL{zk.ACL{Perms: zk.PermAll, Scheme: "world", ID: "anyone"}}
}

func escapePath(path string) string {
	return url.QueryEscape(path)
}

func unescapePath(path string) (string, error) {
	return url.QueryUnescape(path)
}

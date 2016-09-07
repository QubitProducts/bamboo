package service

type Service struct {
	Id string `param:"id"`
	// Acl is present for backwards compatability, and should not be written to.
	// Instead, write to the "Acl" key in Config
	Acl    string            `param:"acl"`
	Config map[string]string `param:"config"`
}

// The storage primitives required by Bamboo from the storage backend
type Storage interface {
	All() ([]Service, error)
	Upsert(service Service) error
	Delete(serviceId string) error
}

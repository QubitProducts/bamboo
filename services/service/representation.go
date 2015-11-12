package service

import (
	"encoding/json"
	"fmt"
)

type ServiceRepr interface {
	Service() Service
	Serialize() ([]byte, error)
}

// ParseServiceRepr attempts parse using each of the ReprParsers in turn. Returns the first
// result to not error, or the error from the last parser
func ParseServiceRepr(body []byte, path string) (repr ServiceRepr, err error) {
	reprParsers := [](func([]byte, string) (ServiceRepr, error)){
		ParseV2ServiceRepr,
		ParseV1ServiceRepr,
	}

	for _, parser := range reprParsers {
		repr, err = parser(body, path)

		if err == nil {
			return
		}
	}
	return nil, err
}

type V1ServiceRepr struct {
	ID  string
	Acl string
}

func ParseV1ServiceRepr(body []byte, path string) (repr ServiceRepr, err error) {
	return &V1ServiceRepr{
		ID:  path,
		Acl: string(body),
	}, nil
}

func (v1 *V1ServiceRepr) Service() Service {
	return Service{
		Id:     v1.ID,
		Acl:    v1.Acl,
		Config: map[string]string{"Acl": v1.Acl},
	}
}

func (v1 *V1ServiceRepr) Serialize() ([]byte, error) {
	return []byte(v1.Acl), nil
}

type V2ServiceRepr struct {
	ID      string            `json:"-"`
	Version string            `json:"version"` // 2 is only valid version for V2ServiceRepr
	Config  map[string]string `json:"config"`
}

func MakeV2ServiceRepr(service Service) *V2ServiceRepr {
	config := make(map[string]string, len(service.Config)+1)
	for k, v := range service.Config {
		config[k] = v
	}
	if service.Acl != "" {
		config["Acl"] = service.Acl
	}
	return NewV2ServiceRepr(service.Id, config)
}

func NewV2ServiceRepr(appID string, config map[string]string) *V2ServiceRepr {
	return &V2ServiceRepr{
		ID:      appID,
		Version: "2",
		Config:  config,
	}
}

func ParseV2ServiceRepr(body []byte, path string) (ServiceRepr, error) {
	var repr V2ServiceRepr
	err := json.Unmarshal(body, &repr)
	if err != nil {
		return nil, err
	}
	if repr.Version != "2" {
		return nil, fmt.Errorf("Service version is not 2 (%s)", repr.Version)
	}
	repr.ID = path

	return &repr, nil
}

func (v2 *V2ServiceRepr) Service() Service {
	acl, _ := v2.Config["Acl"]
	return Service{
		Id:     v2.ID,
		Acl:    acl,
		Config: v2.Config,
	}
}

func (v2 *V2ServiceRepr) Serialize() ([]byte, error) {
	return json.Marshal(&v2)
}

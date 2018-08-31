package registry

import "google.golang.org/grpc/metadata"

type Registry interface {
	String() string
	Register(Service) error
	UnRegister(Service) error
	GetService(string) ([]Service, error)
}

type Service struct {
	Name     string      `json:"name"`
	Version  string      `json:"version"`
	Metadata metadata.MD `json:"metadata"`
	Nodes    []Node      `json:"nodes"`
}

type Node struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

package main

type SeederManifest struct {
	Organizations map[string]*Organization
}

type Organization struct {
	Name   string
	Spaces map[string]*Space
}

type Space struct {
	Name     string
	Apps     map[string]*App
	Services map[string]*Service
	org      *Organization
	seen     []string
}

type App struct {
	Name          string
	Repo          string            `yaml:",omitempty"`
	Path          string            `yaml:",omitempty"`
	Disk          string            `yaml:",omitempty"`
	Memory        string            `yaml:",omitempty"`
	Instances     string            `yaml:",omitempty"`
	Hostname      string            `yaml:",omitempty"`
	Domain        string            `yaml:",omitempty"`
	Buildpack     string            `yaml:",omitempty"`
	Manifest      string            `yaml:",omitempty"`
	ServiceBroker ServiceBroker     `yaml:"service_broker,omitempty"`
	ServiceAccess []Service         `yaml:"service_access,omitempty"`
	Requires      map[string]string `yaml:",omitempty"`
	space         *Space
}

type ServiceBroker struct {
	Name     string
	Username string `yaml:",omitempty"`
	Password string `yaml:",omitempty"`
	Url      string `yaml:",omitempty"`
}

type Service struct {
	Name     string
	Service  string            `yaml:",omitempty"`
	Plan     string            `yaml:",omitempty"`
	Org      string            `yaml:",omitempty"`
	Requires map[string]string `yaml:",omitempty"`
	space    *Space
}

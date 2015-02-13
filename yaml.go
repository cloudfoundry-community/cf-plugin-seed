package main

type SeederManifest struct {
	Organizations []Organization
}

type Organization struct {
	Name   string
	Spaces []Space
}

type Space struct {
	Name     string
	Apps     []App
	Services []Service
}

type App struct {
	Name      string
	Repo      string `yaml:",omitempty"`
	Path      string `yaml:",omitempty"`
	Disk      string `yaml:",omitempty"`
	Memory    string `yaml:",omitempty"`
	Instances string `yaml:",omitempty"`
	Hostname  string `yaml:",omitempty"`
	Domain    string `yaml:",omitempty"`
	Buildpack string `yaml:",omitempty"`
	Manifest  string `yaml:",omitempty"`
}

type Service struct {
	Name    string
	Service string `yaml:",omitempty"`
	Plan    string `yaml:",omitempty"`
}

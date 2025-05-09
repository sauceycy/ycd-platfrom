package model

type Project struct {
	Name        string       `yaml:"name" json:"name"`
	HarborRepo  string       `yaml:"harbor_repo" json:"harbor_repo"`
	Chart       string       `yaml:"chart" json:"chart"`
	Environments []Environment `yaml:"environments" json:"environments"`
}

type Environment struct {
	Name       string   `yaml:"name" json:"name"`
	Cluster    string   `yaml:"cluster" json:"cluster"`
	Namespaces []string `yaml:"namespaces" json:"namespaces"`
}

type Cluster struct {
	Name     string `yaml:"name" json:"name"`
	HelmAPI  string `yaml:"helm_api" json:"helm_api"`
}
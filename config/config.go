package config

import (
	"os"

	"gopkg.in/yaml.v3"
	"ycd-platform/model"
)

type Config struct {
	Projects []model.Project `yaml:"projects"`
	Clusters []model.Cluster `yaml:"clusters"`
}

var Global Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &Global)
}

// SaveConfig 将当前配置保存到文件
func SaveConfig(filePath string) error {
	data, err := yaml.Marshal(Global)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
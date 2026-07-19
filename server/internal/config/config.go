package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// Config 是应用的全部配置，对应 config/config.yaml。
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig 是 HTTP 服务相关配置。
type ServerConfig struct {
	Port int `yaml:"port"`
}

// DatabaseConfig 是数据库相关配置。
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// Load 读取并解析 path 指向的 yaml 配置文件。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	return &cfg, nil
}

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ProviderMaxPrice struct {
	Prompt     float64 `yaml:"prompt" json:"prompt,omitempty"`
	Completion float64 `yaml:"completion" json:"completion,omitempty"`
}

type ProviderConfig struct {
	Order          []string          `yaml:"order" json:"order,omitempty"`
	Allow          []string          `yaml:"allow" json:"allow,omitempty"`
	Deny           []string          `yaml:"deny" json:"deny,omitempty"`
	MaxPrice       *ProviderMaxPrice `yaml:"max_price" json:"max_price,omitempty"`
	AllowFallbacks *bool             `yaml:"allow_fallbacks" json:"allow_fallbacks,omitempty"`
}

type ModelConfig struct {
	ID       string         `yaml:"id"`
	Provider ProviderConfig `yaml:"provider"`
}

type ServerConfig struct {
	Listen string `yaml:"listen"`
}

type RoutingConfig struct {
	ChatPaths []string `yaml:"chat_paths"`
}

type Config struct {
	Server  ServerConfig           `yaml:"server"`
	Models  map[string]ModelConfig `yaml:"models"`
	Routing RoutingConfig          `yaml:"routing"`
}

var C Config

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &C)
}

func ResolveModel(alias string) (ModelConfig, bool) {
	m, ok := C.Models[alias]
	return m, ok
}

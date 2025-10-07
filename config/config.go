package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	BotToken  string         `yaml:"bot_token"`
	ProjectID string         `yaml:"project_id"`
	Device    []Device       `yaml:"devices"`
	Commands  []SlashCommand `yaml:"commands"`
}

type Device struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}

type SlashCommand struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Task        string     `yaml:"task"`
	Variables   []Variable `yaml:"variables,omitempty"`
}

type Variable struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	return &cfg, nil
}

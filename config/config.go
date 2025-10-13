package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type Config struct {
	BotToken  string         `yaml:"bot_token" validate:"required"`
	ProjectID string         `yaml:"project_id" validate:"required"`
	Device    []Device       `yaml:"devices" validate:"required,dive"`
	Commands  []SlashCommand `yaml:"commands" validate:"required,unique_command,dive"`
}

type Device struct {
	Name  string `yaml:"name" validate:"required"`
	Token string `yaml:"token" validate:"required"`
}

type SlashCommand struct {
	Name        string     `yaml:"name" validate:"required,ne=help"`
	Description string     `yaml:"description" validate:"required"`
	Task        string     `yaml:"task" validate:"required"`
	TTL         uint       `yaml:"ttl"`
	Variables   []Variable `yaml:"variables,omitempty" validate:"unique_variable,dive"`
}

type Variable struct {
	Name        string `yaml:"name" validate:"required"`
	Description string `yaml:"description" validate:"required"`
	Type        uint8  `yaml:"type" validate:"gte=3,lte=8|gte=10,lte=11"`
	Required    bool   `yaml:"required"`
}

var validate *validator.Validate

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

	// set up validator
	validate = validator.New(validator.WithRequiredStructEnabled())

	err = validate.RegisterValidation("unique_command", validateUniqueCommandName)
	if err != nil {
		return nil, fmt.Errorf("failed to register unique_command validator: %w", err)
	}

	err = validate.RegisterValidation("unique_variable", validateUniqueVariableName)
	if err != nil {
		return nil, fmt.Errorf("failed to register unique_variable validator: %w", err)
	}

	// validate config
	err = validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate config file: %w", err)
	}

	return &cfg, nil
}

func validateUniqueCommandName(fl validator.FieldLevel) bool {
	commands, ok := fl.Field().Interface().([]SlashCommand)
	if !ok {
		return false
	}

	processedCommands := make(map[string]bool)
	for _, cmd := range commands {
		if _, ok := processedCommands[cmd.Name]; ok {
			return false
		}
		processedCommands[cmd.Name] = true
	}

	return true
}

func validateUniqueVariableName(fl validator.FieldLevel) bool {
	variables, ok := fl.Field().Interface().([]Variable)
	if !ok {
		return false
	}

	processedVariables := make(map[string]bool)
	for _, variable := range variables {
		if _, ok := processedVariables[variable.Name]; ok {
			return false
		}
		processedVariables[variable.Name] = true
	}

	return true
}

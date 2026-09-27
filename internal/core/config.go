// Package core contains the core functionality of CDSL.
package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the CDSL configuration.
type Config map[string]any

// LoadConfig loads the config and shares reference to it
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configFile := filepath.Join(home, ".config", "cdsl", "config.toml")

	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// ComponentExists checks whether the component name passed by user exists or not
func ComponentExists(config *Config, componentName string) (bool, error) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
		return false, nil
	}

	_, exists := components[componentName]

	return exists, nil
}

// ComponentID looks up for the component and returns the configured id
func ComponentID(config *Config, componentName string) (uint8, error) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("no components configured")
	}

	component, ok := components[componentName].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("component not configured: %s", componentName)
	}

	id, ok := component["id"].(int64)
	if !ok {
		return 0, fmt.Errorf("component ID is not configured: %s", componentName)
	}

	if id < 0 || id > 255 {
		return 0, fmt.Errorf("component ID must be between 0 and 255")
	}

	return uint8(id), nil
}

// ValueWithinLimits checks whether the provided float value falls within the configured limits
func ValueWithinLimits(config *Config, componentName string, value float32) (bool, error) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
		return false, nil
	}

	component, ok := components[componentName].(map[string]any)
	if !ok {
		return false, nil
	}

	min, minOK := component["min"].(int64)
	max, maxOK := component["max"].(int64)

	if !minOK || !maxOK {
		return false, nil
	}

	return value >= float32(min) && value <= float32(max), nil
}

// ArduinoPort returns the port of the connected arduino
func ArduinoPort(config *Config) (string, error) {
	arduino, ok := (*config)["arduino"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("arduino serial port is not configured")
	}

	port, ok := arduino["port"].(string)
	if !ok {
		return "", fmt.Errorf("arduino serial port is not configured")
	}

	return port, nil
}

// ArduinoBaudRate returns the configured baud rate
func ArduinoBaudRate(config *Config) (int, error) {
	arduino, ok := (*config)["arduino"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("arduino baud rate is not configured")
	}

	baudRate, ok := arduino["baud_rate"].(int64)
	if !ok {
		return 0, fmt.Errorf("arduino baud rate is not configured")
	}

	return int(baudRate), nil
}

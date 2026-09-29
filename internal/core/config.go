// Package core contains the core functionality of CDSL.
package core

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// ========================================
// Configuration
// ========================================

// Config represents the CDSL configuration.
//
// The configuration is stored as a generic map to keep the rest of the
// application independent of the underlying configuration file format.
type Config map[string]any

// ========================================
// XML Configuration Types
// ========================================

// xmlConfig represents the structure of the CDSL XML configuration file.
type xmlConfig struct {
	XMLName   xml.Name     `xml:"config"`
	Arduino   xmlArduino   `xml:"arduino"`
	Component xmlComponent `xml:"component"`
}

// xmlArduino contains the Arduino serial connection settings.
type xmlArduino struct {
	Port     string `xml:"port"`
	BaudRate int    `xml:"baud_rate"`
}

// xmlComponent contains the configuration for all supported robot components.
type xmlComponent struct {
	Base     xmlComponentConfig `xml:"base"`
	Shoulder xmlComponentConfig `xml:"shoulder"`
	Elbow    xmlComponentConfig `xml:"elbow"`
	Wrist    xmlComponentConfig `xml:"wrist"`
	Hand     xmlComponentConfig `xml:"hand"`
	Claw     xmlComponentConfig `xml:"claw"`
}

// xmlComponentConfig contains the hardware configuration for a single
// robot component.
type xmlComponentConfig struct {
	ID      uint8 `xml:"id"`
	Min     int   `xml:"min"`
	Max     int   `xml:"max"`
	Default int   `xml:"default"`
}

// ========================================
// Configuration Loading
// ========================================

// LoadConfig loads the CDSL configuration from the user's configuration
// directory.
//
// The configuration file is expected at:
//
//	~/.config/cdsl/config.xml
//
// XML parsing is isolated to this function so the rest of the application
// operates on the generic Config representation.
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configFile := filepath.Join(home, ".config", "cdsl", "config.xml")

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var xmlConfig xmlConfig

	if err := xml.Unmarshal(data, &xmlConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Convert the XML representation into the generic configuration structure
	// used by the rest of the application.
	config := Config{
		"arduino": map[string]any{
			"port":      xmlConfig.Arduino.Port,
			"baud_rate": int64(xmlConfig.Arduino.BaudRate),
		},
		"component": map[string]any{
			"base": map[string]any{
				"id":      int64(xmlConfig.Component.Base.ID),
				"min":     int64(xmlConfig.Component.Base.Min),
				"max":     int64(xmlConfig.Component.Base.Max),
				"default": int64(xmlConfig.Component.Base.Default),
			},
			"shoulder": map[string]any{
				"id":      int64(xmlConfig.Component.Shoulder.ID),
				"min":     int64(xmlConfig.Component.Shoulder.Min),
				"max":     int64(xmlConfig.Component.Shoulder.Max),
				"default": int64(xmlConfig.Component.Shoulder.Default),
			},
			"elbow": map[string]any{
				"id":      int64(xmlConfig.Component.Elbow.ID),
				"min":     int64(xmlConfig.Component.Elbow.Min),
				"max":     int64(xmlConfig.Component.Elbow.Max),
				"default": int64(xmlConfig.Component.Elbow.Default),
			},
			"wrist": map[string]any{
				"id":      int64(xmlConfig.Component.Wrist.ID),
				"min":     int64(xmlConfig.Component.Wrist.Min),
				"max":     int64(xmlConfig.Component.Wrist.Max),
				"default": int64(xmlConfig.Component.Wrist.Default),
			},
			"hand": map[string]any{
				"id":      int64(xmlConfig.Component.Hand.ID),
				"min":     int64(xmlConfig.Component.Hand.Min),
				"max":     int64(xmlConfig.Component.Hand.Max),
				"default": int64(xmlConfig.Component.Hand.Default),
			},
			"claw": map[string]any{
				"id":      int64(xmlConfig.Component.Claw.ID),
				"min":     int64(xmlConfig.Component.Claw.Min),
				"max":     int64(xmlConfig.Component.Claw.Max),
				"default": int64(xmlConfig.Component.Claw.Default),
			},
		},
	}

	return &config, nil
}

// ========================================
// Component Operations
// ========================================

// ComponentExists reports whether the specified component exists in the
// configuration.
func ComponentExists(config *Config, componentName string) (bool, error) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
		return false, nil
	}

	_, exists := components[componentName]

	return exists, nil
}

// ComponentID returns the configured ID of the specified component.
//
// The ID must be within the range of an unsigned 8-bit integer.
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

// ValueWithinLimits reports whether value falls within the configured
// minimum and maximum values for the specified component.
//
// The configured limits are inclusive.
func ValueWithinLimits(config *Config, componentName string, value float32) (bool, error) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
		return false, nil
	}

	component, ok := components[componentName].(map[string]any)
	if !ok {
		return false, nil
	}

	minValue, minOK := component["min"].(int64)
	maxValue, maxOK := component["max"].(int64)

	if !minOK || !maxOK {
		return false, nil
	}

	return value >= float32(minValue) && value <= float32(maxValue), nil
}

// ComponentDefault returns the configured default value for the specified
// component.
func ComponentDefault(config *Config, componentName string) (float32, error) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("no components configured")
	}

	component, ok := components[componentName].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("component not configured: %s", componentName)
	}

	defaultValue, ok := component["default"].(int64)
	if !ok {
		return 0, fmt.Errorf("component default is not configured: %s", componentName)
	}

	return float32(defaultValue), nil
}

// ========================================
// Arduino Configuration
// ========================================

// ArduinoPort returns the serial port configured for the Arduino.
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

// ArduinoBaudRate returns the configured baud rate for the Arduino.
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

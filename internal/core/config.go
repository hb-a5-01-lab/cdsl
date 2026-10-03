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
	XMLName   xml.Name      `xml:"config"`
	Arduino   xmlArduino    `xml:"arduino"`
	Component xmlComponents `xml:"component"`
}

// xmlArduino contains the Arduino serial connection settings.
type xmlArduino struct {
	Port     string `xml:"port"`
	BaudRate int    `xml:"baud_rate"`
}

// xmlComponents represents the collection of robot components.
//
// Components are intentionally stored dynamically rather than as named
// struct fields. This allows new components to be added to the XML
// configuration without requiring changes to the Go source code.
type xmlComponents struct {
	Elements []xmlComponent `xml:",any"`
}

// xmlComponent represents a single robot component.
//
// The component name is obtained from the XML element itself, while its
// configuration is parsed from the element's child elements.
type xmlComponent struct {
	XMLName xml.Name `xml:""`

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
		"component": make(map[string]any),
	}

	components := config["component"].(map[string]any)

	for _, component := range xmlConfig.Component.Elements {
		name := component.XMLName.Local

		components[name] = map[string]any{
			"id":      int64(component.ID),
			"min":     int64(component.Min),
			"max":     int64(component.Max),
			"default": int64(component.Default),
		}
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

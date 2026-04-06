package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFile = ".figmarc.json"

// Config holds project-level defaults persisted to .figmarc.json.
type Config struct {
	FileKey   string `json:"fileKey,omitempty"`
	Theme     string `json:"theme,omitempty"`
	Page      int    `json:"page,omitempty"`
	ExportDir string `json:"exportDir,omitempty"`
}

// Load reads .figmarc.json from the current directory or parents.
// Returns a zero Config (with defaults) if no file is found.
func Load() (*Config, error) {
	path, err := findConfigFile()
	if err != nil {
		return &Config{Theme: "default"}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if c.Theme == "" {
		c.Theme = "default"
	}
	return &c, nil
}

// Save writes the config to .figmarc.json in the current directory.
func Save(c *Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile(configFile, data, 0644)
}

// Init creates a new .figmarc.json with the given name as a starting point.
func Init(name string) error {
	c := &Config{
		Theme: "default",
	}
	return Save(c)
}

// Get returns the value of a config key.
func Get(key string) (string, error) {
	c, err := Load()
	if err != nil {
		return "", err
	}
	switch key {
	case "fileKey":
		return c.FileKey, nil
	case "theme":
		return c.Theme, nil
	case "page":
		return fmt.Sprintf("%d", c.Page), nil
	case "exportDir":
		return c.ExportDir, nil
	default:
		return "", fmt.Errorf("unknown config key %q", key)
	}
}

// Set updates a config key.
func Set(key, value string) error {
	c, err := Load()
	if err != nil {
		return err
	}
	switch key {
	case "fileKey":
		c.FileKey = value
	case "theme":
		c.Theme = value
	case "page":
		var page int
		if _, err := fmt.Sscanf(value, "%d", &page); err != nil {
			return fmt.Errorf("page must be an integer: %w", err)
		}
		c.Page = page
	case "exportDir":
		c.ExportDir = value
	default:
		return fmt.Errorf("unknown config key %q (valid: fileKey, theme, page, exportDir)", key)
	}
	return Save(c)
}

func findConfigFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		p := filepath.Join(dir, configFile)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found", configFile)
		}
		dir = parent
	}
}

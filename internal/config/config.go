// Package config provides configuration loading and management for grule-lint.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	"gopkg.in/yaml.v3"

	"github.com/adarshjos/grule-lint/internal/diagnostic"
)

const ConfigFileName = ".grl-lint.yaml"

type Config struct {
	Rules      map[string]string `yaml:"rules"`
	Exclude    []string          `yaml:"exclude"`
	Include    []string          `yaml:"include"`
	Complexity ComplexityConfig  `yaml:"complexity"`
	Naming     NamingConfig      `yaml:"naming"`
}

type ComplexityConfig struct {
	MaxConditions int `yaml:"max_conditions"`
}

// NamingConfig holds configuration for the naming convention rule.
type NamingConfig struct {
	// Convention specifies the naming convention to enforce.
	// Valid values: "PascalCase", "camelCase", "snake_case", "kebab-case"
	// Default: "PascalCase"
	Convention string `yaml:"convention"`
}

func DefaultConfig() *Config {
	return &Config{
		Rules:      make(map[string]string),
		Exclude:    []string{},
		Include:    []string{"**/*.grl"},
		Complexity: ComplexityConfig{MaxConditions: 5},
		Naming:     NamingConfig{Convention: "PascalCase"},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}
	return config, nil
}

func LoadFromDirectory(dir string) (*Config, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving directory path %s: %w", dir, err)
	}

	current := absDir
	for {
		configPath := filepath.Join(current, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return Load(configPath)
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return DefaultConfig(), nil
}

func (c *Config) GetRuleSeverity(ruleID string, defaultSeverity diagnostic.Severity) *diagnostic.Severity {
	if severityStr, ok := c.Rules[ruleID]; ok {
		switch severityStr {
		case "off", "disabled":
			return nil
		case "error":
			s := diagnostic.SeverityError
			return &s
		case "warning", "warn":
			s := diagnostic.SeverityWarning
			return &s
		case "info":
			s := diagnostic.SeverityInfo
			return &s
		case "hint":
			s := diagnostic.SeverityHint
			return &s
		}
	}
	return &defaultSeverity
}

func (c *Config) IsRuleEnabled(ruleID string) bool {
	if severityStr, ok := c.Rules[ruleID]; ok {
		return severityStr != "off" && severityStr != "disabled"
	}
	return true
}

func (c *Config) ShouldExclude(file string) bool {
	// Normalize path separators for cross-platform matching
	normalizedFile := filepath.ToSlash(file)

	for _, pattern := range c.Exclude {
		// Use doublestar for proper glob matching with ** support
		matched, err := doublestar.Match(pattern, normalizedFile)
		if err == nil && matched {
			return true
		}

		// Also try matching against just the basename
		matched, err = doublestar.Match(pattern, filepath.Base(file))
		if err == nil && matched {
			return true
		}
	}
	return false
}

func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	for k, v := range other.Rules {
		c.Rules[k] = v
	}

	if len(other.Exclude) > 0 {
		c.Exclude = append(c.Exclude, other.Exclude...)
	}

	if len(other.Include) > 0 {
		c.Include = other.Include
	}

	if other.Complexity.MaxConditions > 0 {
		c.Complexity.MaxConditions = other.Complexity.MaxConditions
	}

	if other.Naming.Convention != "" {
		c.Naming.Convention = other.Naming.Convention
	}
}

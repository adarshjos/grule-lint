package lint

import (
	"github.com/adarshjos/grule-lint/internal/config"
)

// ConfigFileName is the default configuration file name.
const ConfigFileName = config.ConfigFileName

// Config holds linting configuration options.
type Config struct {
	c *config.Config
}

// DefaultConfig returns a new Config with default settings.
func DefaultConfig() *Config {
	return &Config{c: config.DefaultConfig()}
}

// LoadConfig loads configuration from the specified file path.
func LoadConfig(path string) (*Config, error) {
	c, err := config.Load(path)
	if err != nil {
		return nil, err
	}
	return &Config{c: c}, nil
}

// LoadConfigFromDirectory searches for a .grl-lint.yaml file starting
// from the given directory and walking up to parent directories.
func LoadConfigFromDirectory(dir string) (*Config, error) {
	c, err := config.LoadFromDirectory(dir)
	if err != nil {
		return nil, err
	}
	return &Config{c: c}, nil
}

// Rules returns the rule severity overrides.
func (c *Config) Rules() map[string]string {
	return c.c.Rules
}

// SetRule sets the severity for a specific rule.
// Valid severities: "error", "warning", "info", "hint", "off"
func (c *Config) SetRule(ruleID, severity string) {
	c.c.Rules[ruleID] = severity
}

// Exclude returns the list of exclude patterns.
func (c *Config) Exclude() []string {
	return c.c.Exclude
}

// SetExclude sets the exclude patterns.
func (c *Config) SetExclude(patterns []string) {
	c.c.Exclude = patterns
}

// AddExclude adds an exclude pattern.
func (c *Config) AddExclude(pattern string) {
	c.c.Exclude = append(c.c.Exclude, pattern)
}

// Include returns the list of include patterns.
func (c *Config) Include() []string {
	return c.c.Include
}

// SetInclude sets the include patterns.
func (c *Config) SetInclude(patterns []string) {
	c.c.Include = patterns
}

// MaxConditions returns the maximum number of conditions allowed
// before the high-complexity rule triggers.
func (c *Config) MaxConditions() int {
	return c.c.Complexity.MaxConditions
}

// SetMaxConditions sets the maximum number of conditions allowed.
func (c *Config) SetMaxConditions(max int) {
	c.c.Complexity.MaxConditions = max
}

// NamingConvention returns the naming convention for rule names.
// Valid values: "PascalCase", "camelCase", "snake_case", "kebab-case"
func (c *Config) NamingConvention() string {
	return c.c.Naming.Convention
}

// SetNamingConvention sets the naming convention for rule names.
// Valid values: "PascalCase", "camelCase", "snake_case", "kebab-case"
func (c *Config) SetNamingConvention(convention string) {
	c.c.Naming.Convention = convention
}

// IsRuleEnabled returns whether a rule is enabled.
func (c *Config) IsRuleEnabled(ruleID string) bool {
	return c.c.IsRuleEnabled(ruleID)
}

// ShouldExclude returns whether a file should be excluded from linting.
func (c *Config) ShouldExclude(file string) bool {
	return c.c.ShouldExclude(file)
}

// Merge merges another config into this one.
// Values from other take precedence.
func (c *Config) Merge(other *Config) {
	if other != nil {
		c.c.Merge(other.c)
	}
}

// internal returns the underlying internal config.
func (c *Config) internal() *config.Config {
	return c.c
}

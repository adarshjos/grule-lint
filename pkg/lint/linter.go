package lint

import (
	"github.com/adarshjos/grule-lint/internal/linter"
	"github.com/adarshjos/grule-lint/internal/rules"
)

// Linter provides GRL file linting capabilities.
type Linter struct {
	l      *linter.Linter
	config *Config
}

// New creates a new Linter with default configuration.
func New() *Linter {
	return &Linter{
		l:      linter.New(),
		config: DefaultConfig(),
	}
}

// NewWithConfig creates a new Linter with the specified configuration.
func NewWithConfig(cfg *Config) *Linter {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Create registry config from the public config
	regCfg := rules.RegistryConfig{
		NamingConvention: cfg.NamingConvention(),
		MaxConditions:    cfg.MaxConditions(),
	}

	return &Linter{
		l:      linter.NewWithConfig(regCfg),
		config: cfg,
	}
}

// Config returns the linter's configuration.
func (l *Linter) Config() *Config {
	return l.config
}

// LintFile lints a single GRL file and returns the results.
func (l *Linter) LintFile(file string) (*Result, error) {
	if l.config.ShouldExclude(file) {
		return &Result{ds: nil}, nil
	}

	ds, err := l.l.LintFile(file)
	if err != nil {
		return nil, err
	}
	return wrapDiagnosticSet(ds), nil
}

// LintString lints GRL content from a string.
// The file parameter is used for diagnostic reporting.
func (l *Linter) LintString(file, content string) *Result {
	ds := l.l.LintString(file, content)
	return wrapDiagnosticSet(ds)
}

// LintFiles lints multiple GRL files.
func (l *Linter) LintFiles(files []string) (*Result, error) {
	// Filter excluded files
	var filtered []string
	for _, f := range files {
		if !l.config.ShouldExclude(f) {
			filtered = append(filtered, f)
		}
	}

	ds, err := l.l.LintFiles(filtered)
	if err != nil {
		return nil, err
	}
	return wrapDiagnosticSet(ds), nil
}

// LintDirectory lints all GRL files in a directory recursively.
func (l *Linter) LintDirectory(dir string) (*Result, error) {
	ds, err := l.l.LintDirectory(dir)
	if err != nil {
		return nil, err
	}
	return wrapDiagnosticSet(ds), nil
}

// LintPaths lints files and/or directories.
func (l *Linter) LintPaths(paths []string) (*Result, error) {
	ds, err := l.l.LintPaths(paths)
	if err != nil {
		return nil, err
	}
	return wrapDiagnosticSet(ds), nil
}

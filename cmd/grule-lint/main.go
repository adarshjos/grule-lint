package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/adarshjos/grule-lint/internal/config"
	"github.com/adarshjos/grule-lint/internal/diagnostic"
	"github.com/adarshjos/grule-lint/internal/linter"
	"github.com/adarshjos/grule-lint/internal/reporter"
)

var (
	// Set via ldflags
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"

	// CLI flags
	configFlag  string
	outputFlag  string
	ruleFlags   []string
	excludeFlag []string
	quietFlag   bool
	noColorFlag bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "grule-lint [files/directories...]",
		Short: "A linter for GRL (Grule Rule Language) files",
		Long: `grule-lint is a static analysis tool for GRL files used with the Grule Rule Engine.

It checks for syntax errors and common issues in your rule definitions.

Examples:
  grule-lint rules.grl
  grule-lint rules/
  grule-lint rules/*.grl
  grule-lint --config .grl-lint.yaml rules/
  grule-lint --quiet rules/`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildTime),
		Args:    cobra.MinimumNArgs(1),
		RunE:    runLint,
	}

	// Add flags
	rootCmd.Flags().StringVarP(&configFlag, "config", "c", "", "Path to config file (default: .grl-lint.yaml)")
	rootCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().StringArrayVarP(&ruleFlags, "rule", "r", nil, "Enable only specific rules (can be repeated)")
	rootCmd.Flags().StringArrayVarP(&excludeFlag, "exclude", "e", nil, "Exclude file patterns (can be repeated)")
	rootCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Only show errors, not warnings/info")
	rootCmd.Flags().BoolVar(&noColorFlag, "no-color", false, "Disable colored output")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runLint(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := loadConfig(args)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Apply CLI overrides to config
	applyCliOverrides(cfg)

	// Create linter
	l := linter.New()

	// Filter paths based on exclusions
	paths := filterPaths(args, cfg)

	if len(paths) == 0 {
		fmt.Println("No files to lint after applying exclusions.")
		return nil
	}

	// Lint all provided paths
	diagnostics, err := l.LintPaths(paths)
	if err != nil {
		return fmt.Errorf("linting failed: %w", err)
	}

	// Filter diagnostics based on config
	filtered := filterDiagnostics(diagnostics.All(), cfg)

	// Set up output writer
	var output io.Writer = os.Stdout
	if outputFlag != "" {
		file, err := os.Create(outputFlag)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer func() { _ = file.Close() }()
		output = file
	}

	// Create text reporter
	useColors := !noColorFlag && outputFlag == "" && isTerminal()
	rep := reporter.NewTextReporter(output, useColors)

	// Create a new DiagnosticSet from filtered diagnostics
	ds := diagnostic.NewDiagnosticSet()
	ds.AddAll(filtered)

	if err := rep.Report(ds); err != nil {
		return fmt.Errorf("reporting failed: %w", err)
	}

	if ds.Count() == 0 {
		if _, err := fmt.Fprintln(output, "No issues found."); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
	}

	// Exit with error code if there are errors
	if ds.HasErrors() {
		os.Exit(1)
	}

	return nil
}

// loadConfig loads configuration from file or defaults.
func loadConfig(args []string) (*config.Config, error) {
	// If explicit config path provided, load it
	if configFlag != "" {
		return config.Load(configFlag)
	}

	// Try to find config by walking up from the first path
	if len(args) > 0 {
		startDir := args[0]
		info, err := os.Stat(startDir)
		if err == nil {
			if !info.IsDir() {
				startDir = filepath.Dir(startDir)
			}
			return config.LoadFromDirectory(startDir)
		}
	}

	// Return defaults
	return config.DefaultConfig(), nil
}

// applyCliOverrides applies CLI flags to the config.
func applyCliOverrides(cfg *config.Config) {
	// Add CLI exclude patterns
	if len(excludeFlag) > 0 {
		cfg.Exclude = append(cfg.Exclude, excludeFlag...)
	}

	// If specific rules are enabled via CLI, enable only those rules
	if len(ruleFlags) > 0 {
		// Build a set of enabled rule IDs for O(1) lookup
		enabledRules := make(map[string]bool, len(ruleFlags))
		for _, ruleID := range ruleFlags {
			enabledRules[ruleID] = true
		}

		// Disable all rules that are not in the enabled set
		allRuleIDs := []string{
			"GRL001", "GRL002", "GRL003", "GRL004", "GRL005", "GRL006",
			"GRL007", "GRL008", "GRL009", "GRL010", "GRL011", "GRL012",
		}
		for _, ruleID := range allRuleIDs {
			if !enabledRules[ruleID] {
				cfg.Rules[ruleID] = "off"
			}
		}
	}
}

// filterPaths filters file paths based on exclusions.
func filterPaths(args []string, cfg *config.Config) []string {
	var result []string

	for _, path := range args {
		info, err := os.Stat(path)
		if err != nil {
			// Keep invalid paths - let the linter handle the error
			result = append(result, path)
			continue
		}

		if info.IsDir() {
			// For directories, we'll let the linter walk them
			// The linter will apply exclusions per-file
			result = append(result, path)
		} else {
			// Check if file should be excluded
			if !cfg.ShouldExclude(path) {
				result = append(result, path)
			}
		}
	}

	return result
}

// filterDiagnostics filters diagnostics based on config and quiet mode.
func filterDiagnostics(diags []diagnostic.Diagnostic, cfg *config.Config) []diagnostic.Diagnostic {
	var result []diagnostic.Diagnostic

	for _, d := range diags {
		// Check if rule is enabled
		if !cfg.IsRuleEnabled(d.RuleID) {
			continue
		}

		// Apply severity override
		severity := cfg.GetRuleSeverity(d.RuleID, d.Severity)
		if severity == nil {
			continue // Rule disabled
		}

		// Apply quiet mode filter
		if quietFlag {
			if *severity != diagnostic.SeverityError {
				continue
			}
		}

		// Update severity if overridden
		d.Severity = *severity

		result = append(result, d)
	}

	return result
}

// isTerminal checks if stdout is a terminal.
func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

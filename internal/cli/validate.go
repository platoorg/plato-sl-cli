package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/spf13/cobra"
	"github.com/platoorg/platosl-cli/internal/config"
	platoCue "github.com/platoorg/platosl-cli/internal/cue"
	platoErrors "github.com/platoorg/platosl-cli/internal/errors"
)

var (
	validateStrict bool
)

var validateCmd = &cobra.Command{
	Use:   "validate [file or directory]",
	Short: "Validate CUE schemas",
	Long: `Validate CUE schemas for correctness and completeness.

If a file or directory is specified, validates only that path.
Otherwise, validates all schema paths from platosl.yaml.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "strict validation (requires all fields to be concrete)")
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Determine what to validate
	var paths []string
	useConfig := false

	if len(args) > 0 {
		// Validate specific path
		path := args[0]

		// Check if path exists (use absolute for stat check)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}

		// Keep relative path for CUE loader (it doesn't like absolute paths)
		paths = []string{path}
		PrintVerbose("Validating: %s", path)
	} else {
		// Load config and validate configured paths
		useConfig = true
		cfg, err := config.Load(GetConfigFile())
		if err != nil {
			return err
		}

		// Override strict setting if specified on command line
		strict := cfg.Validation.Strict
		if validateStrict {
			strict = true
		}
		validateStrict = strict

		// Collect all schema paths (keep relative for CUE)
		for _, schemaPath := range cfg.Schemas {
			// Validate path exists
			if _, err := os.Stat(schemaPath); err != nil {
				PrintError("Failed to access path %s: %v", schemaPath, err)
				continue
			}
			paths = append(paths, schemaPath)
		}

		if len(paths) == 0 {
			return fmt.Errorf("no schema paths configured in platosl.yaml")
		}

		PrintVerbose("Validating %d schema path(s) from config", len(paths))
	}

	// Create loader and validator
	loader := platoCue.NewLoader()
	validator := platoCue.NewValidator(validateStrict)

	// Track validation results
	var allErrors []*platoErrors.Error
	validatedFiles := 0

	// Expand directories to find all CUE packages
	var expandedPaths []string
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			allErrors = append(allErrors, platoErrors.Newf(
				platoErrors.ErrorTypeFileSystem,
				"cannot access %s: %v", path, err,
			))
			continue
		}

		if info.IsDir() {
			// Find all subdirectories with CUE files
			subPaths, err := findCuePackages(path)
			if err != nil {
				allErrors = append(allErrors, platoErrors.Newf(
					platoErrors.ErrorTypeFileSystem,
					"failed to scan directory %s: %v", path, err,
				))
				continue
			}
			if len(subPaths) == 0 {
				// No CUE files found in subdirectories, try the directory itself
				expandedPaths = append(expandedPaths, path)
			} else {
				expandedPaths = append(expandedPaths, subPaths...)
			}
		} else {
			expandedPaths = append(expandedPaths, path)
		}
	}

	// Validate each path
	for _, path := range expandedPaths {
		info, err := os.Stat(path)
		if err != nil {
			allErrors = append(allErrors, platoErrors.Newf(
				platoErrors.ErrorTypeFileSystem,
				"cannot access %s: %v", path, err,
			))
			continue
		}

		var val cue.Value
		if info.IsDir() {
			PrintVerbose("Loading directory: %s", path)
			val, err = loader.LoadDir(path)
		} else {
			PrintVerbose("Loading file: %s", filepath.Base(path))
			val, err = loader.LoadFile(path)
		}

		if err != nil {
			allErrors = append(allErrors, platoErrors.Wrapf(
				platoErrors.ErrorTypeValidation,
				err,
				"failed to load %s", path,
			))
			continue
		}

		// Validate
		result := validator.Validate(val)
		validatedFiles++

		if !result.Valid {
			for _, verr := range result.Errors {
				allErrors = append(allErrors, platoErrors.New(
					platoErrors.ErrorTypeValidation,
					verr.Message,
				).WithLocation(verr.File, verr.Line, verr.Column).WithSuggestion(verr.Suggestion))
			}
		}
	}

	// Report results
	if len(allErrors) > 0 {
		PrintError("Validation failed\n")
		for _, err := range allErrors {
			fmt.Fprintln(os.Stderr, err.Format())
			fmt.Fprintln(os.Stderr)
		}
		return fmt.Errorf("found %d error(s)", len(allErrors))
	}

	// Success
	if useConfig {
		PrintSuccess("All schemas valid (%d path(s) checked)", len(paths))
	} else {
		PrintSuccess("Schema valid")
	}

	return nil
}

// findCuePackages finds all directories containing CUE files recursively
func findCuePackages(rootPath string) ([]string, error) {
	var packages []string
	seen := make(map[string]bool)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and cue.mod
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "cue.mod") {
			return filepath.SkipDir
		}

		// If it's a .cue file, add its directory
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".cue") {
			dir := filepath.Dir(path)
			if !seen[dir] {
				seen[dir] = true
				packages = append(packages, dir)
			}
		}

		return nil
	})

	return packages, err
}

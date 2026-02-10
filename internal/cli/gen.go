package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/spf13/cobra"
	"github.com/platoorg/plato-sl-cli/internal/config"
	platoCue "github.com/platoorg/plato-sl-cli/internal/cue"
	"github.com/platoorg/plato-sl-cli/internal/errors"
	"github.com/platoorg/plato-sl-cli/internal/generator"

	// Import generators to register them
	_ "github.com/platoorg/plato-sl-cli/internal/generator/elixir"
	_ "github.com/platoorg/plato-sl-cli/internal/generator/golang"
	_ "github.com/platoorg/plato-sl-cli/internal/generator/jsonschema"
	_ "github.com/platoorg/plato-sl-cli/internal/generator/typescript"
	_ "github.com/platoorg/plato-sl-cli/internal/generator/zod"
)

var (
	genOutput string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code from CUE schemas",
	Long: `Generate code from CUE schemas to various target languages.

Available generators:
  typescript  - Generate TypeScript interfaces
  zod         - Generate Zod schemas with inferred TypeScript types
  jsonschema  - Generate JSON Schema
  go          - Generate Go structs
  elixir      - Generate Elixir typespecs`,
}

var genTypescriptCmd = &cobra.Command{
	Use:   "typescript",
	Short: "Generate TypeScript interfaces",
	Long: `Generate TypeScript interfaces from CUE definitions.

By default, generates to the output specified in platosl.yaml.
Use --output to override.`,
	RunE: runGenTypescript,
}

var genJsonSchemaCmd = &cobra.Command{
	Use:   "jsonschema",
	Short: "Generate JSON Schema",
	Long:  `Generate JSON Schema (draft 2020-12) from CUE definitions.`,
	RunE:  runGenJsonSchema,
}

var genGoCmd = &cobra.Command{
	Use:   "go",
	Short: "Generate Go structs",
	Long:  `Generate Go struct types with JSON tags from CUE definitions.`,
	RunE:  runGenGo,
}

var genElixirCmd = &cobra.Command{
	Use:   "elixir",
	Short: "Generate Elixir typespecs",
	Long:  `Generate Elixir typespecs and structs from CUE definitions.`,
	RunE:  runGenElixir,
}

var genZodCmd = &cobra.Command{
	Use:   "zod",
	Short: "Generate Zod schemas with TypeScript types",
	Long:  `Generate Zod validation schemas with inferred TypeScript types from CUE definitions.`,
	RunE:  runGenZod,
}

var (
	genGoPackage     string
	genElixirModule  string
)

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(genTypescriptCmd)
	genCmd.AddCommand(genJsonSchemaCmd)
	genCmd.AddCommand(genGoCmd)
	genCmd.AddCommand(genElixirCmd)
	genCmd.AddCommand(genZodCmd)

	// TypeScript flags
	genTypescriptCmd.Flags().StringVarP(&genOutput, "output", "o", "", "output file path")

	// JSON Schema flags
	genJsonSchemaCmd.Flags().StringVarP(&genOutput, "output", "o", "", "output file path")

	// Go flags
	genGoCmd.Flags().StringVarP(&genOutput, "output", "o", "", "output file path")
	genGoCmd.Flags().StringVar(&genGoPackage, "package", "", "Go package name")

	// Elixir flags
	genElixirCmd.Flags().StringVarP(&genOutput, "output", "o", "", "output file path")
	genElixirCmd.Flags().StringVar(&genElixirModule, "module", "", "Elixir module name")

	// Zod flags
	genZodCmd.Flags().StringVarP(&genOutput, "output", "o", "", "output file path")
}

func runGenTypescript(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load(GetConfigFile())
	if err != nil {
		return err
	}

	// Get TypeScript generator config
	genCfg, ok := cfg.Generate["typescript"]
	if !ok {
		genCfg = config.GenConfig{
			Enabled: true,
			Output:  "generated/types.ts",
			Options: make(map[string]interface{}),
		}
	}

	// Override output if specified
	if genOutput != "" {
		genCfg.Output = genOutput
	}

	// Validate output path
	if genCfg.Output == "" {
		return fmt.Errorf("no output file specified (use --output or configure in platosl.yaml)")
	}

	PrintVerbose("Generating TypeScript to: %s", genCfg.Output)

	// Load and validate schemas
	val, err := loadAndValidateSchemas(cfg, "TypeScript")
	if err != nil {
		return err
	}

	// Get generator
	gen, err := generator.Get("typescript")
	if err != nil {
		e := errors.Wrap(errors.ErrorTypeInternal, err, "TypeScript generator not registered")
		e = e.WithSuggestion("This is an internal error. Please report this issue")
		PrintError(e.Format())
		return e
	}

	// Create generator context
	ctx := generator.NewContext(val, cfg, genCfg)

	// Validate
	PrintVerbose("Validating generator requirements")
	if err := gen.Validate(ctx); err != nil {
		e := errors.Wrap(errors.ErrorTypeValidation, err, "generator validation failed")
		e = e.WithSuggestion("The schema structure may not be compatible with TypeScript generation")
		PrintError(e.Format())
		return e
	}

	// Generate
	PrintVerbose("Generating TypeScript code")
	output, err := gen.Generate(ctx)
	if err != nil {
		e := errors.Wrap(errors.ErrorTypeGeneration, err, "TypeScript generation failed")
		e = e.WithSuggestion("Check that your schema definitions are valid and exportable")
		PrintError(e.Format())
		return e
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(genCfg.Output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		e := errors.Wrap(errors.ErrorTypeFileSystem, err, fmt.Sprintf("failed to create output directory: %s", outputDir))
		e = e.WithSuggestion("Check that you have write permissions for the output directory")
		PrintError(e.Format())
		return e
	}

	// Write output
	if err := os.WriteFile(genCfg.Output, output, 0644); err != nil {
		e := errors.Wrap(errors.ErrorTypeFileSystem, err, fmt.Sprintf("failed to write output file: %s", genCfg.Output))
		e = e.WithSuggestion("Check that you have write permissions for the output file")
		PrintError(e.Format())
		return e
	}

	// Success
	stats := fmt.Sprintf("%d bytes", len(output))
	PrintSuccess("Generated TypeScript: %s (%s)", filepath.Base(genCfg.Output), stats)

	return nil
}

func runGenZod(cmd *cobra.Command, args []string) error {
	return runGenerator("zod", map[string]interface{}{})
}

// runGenAll generates all enabled generators
func runGenAll(cfg *config.Config) error {
	var generated []string
	var genErrors []string

	// Load and validate schemas once for all generators
	PrintVerbose("Loading and validating schemas for all generators")
	loader := platoCue.NewLoader()
	var allPaths []string
	for _, schemaPath := range cfg.Schemas {
		absPath, err := filepath.Abs(schemaPath)
		if err != nil {
			e := errors.Wrap(errors.ErrorTypeFileSystem, err, fmt.Sprintf("failed to resolve schema path: %s", schemaPath))
			PrintError(e.Format())
			return e
		}
		allPaths = append(allPaths, absPath)
	}

	if len(allPaths) == 0 {
		e := errors.New(errors.ErrorTypeConfig, "no schema paths configured")
		e = e.WithSuggestion("Add schema directories to the 'schemas' section in platosl.yaml")
		PrintError(e.Format())
		return e
	}

	val, err := loader.LoadPaths(allPaths)
	if err != nil {
		// Provide context-specific suggestions
		suggestion := "Check your CUE files for syntax errors. Run 'cue vet' directly for more details"
		if strings.Contains(err.Error(), "cannot use absolute directory") {
			suggestion = "CUE module configuration issue. Try using relative paths in platosl.yaml or ensure you have a cue.mod directory"
		} else if strings.Contains(err.Error(), "import failed") {
			suggestion = "Check that all imported packages are available in your cue.mod directory"
		} else if strings.Contains(err.Error(), "cannot find package") {
			suggestion = "Verify that the schema paths in platosl.yaml point to valid CUE packages"
		}

		e := errors.Wrap(errors.ErrorTypeValidation, err, "failed to load schemas")
		e = e.WithSuggestion(suggestion)
		PrintError(e.Format())
		return e
	}

	// Validate schemas once
	validationErrors := validateSchemas(val, "all generators")
	if len(validationErrors) > 0 {
		PrintError("Schema validation failed with %d error(s):\n", len(validationErrors))
		for _, err := range validationErrors {
			PrintError(err.Format())
			fmt.Fprintln(os.Stderr)
		}
		return fmt.Errorf("schema validation failed")
	}

	// Generate for each enabled generator
	for name, genCfg := range cfg.Generate {
		if !genCfg.Enabled {
			PrintVerbose("Skipping disabled generator: %s", name)
			continue
		}

		PrintInfo("Generating %s...", name)

		// Get generator
		gen, err := generator.Get(name)
		if err != nil {
			genErrors = append(genErrors, fmt.Sprintf("%s: generator not registered", name))
			continue
		}

		// Create context and generate
		ctx := generator.NewContext(val, cfg, genCfg)

		if err := gen.Validate(ctx); err != nil {
			genErrors = append(genErrors, fmt.Sprintf("%s: validation failed: %v", name, err))
			continue
		}

		output, err := gen.Generate(ctx)
		if err != nil {
			genErrors = append(genErrors, fmt.Sprintf("%s: generation failed: %v", name, err))
			continue
		}

		// Write output
		outputDir := filepath.Dir(genCfg.Output)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			genErrors = append(genErrors, fmt.Sprintf("%s: failed to create output directory: %s", name, outputDir))
			continue
		}

		if err := os.WriteFile(genCfg.Output, output, 0644); err != nil {
			genErrors = append(genErrors, fmt.Sprintf("%s: failed to write output file: %s", name, genCfg.Output))
			continue
		}

		generated = append(generated, name)
		PrintSuccess("  ✓ %s: %s", name, genCfg.Output)
	}

	// Report results
	if len(genErrors) > 0 {
		PrintError("\nGeneration completed with errors:")
		for _, e := range genErrors {
			PrintError("  %s", e)
		}
	}

	if len(generated) > 0 {
		fmt.Println()
		PrintSuccess("Generated %d target(s): %s", len(generated), strings.Join(generated, ", "))
	}

	if len(genErrors) > 0 {
		return fmt.Errorf("generation completed with %d error(s)", len(genErrors))
	}

	return nil
}

func runGenJsonSchema(cmd *cobra.Command, args []string) error {
	return runGenerator("jsonschema", map[string]interface{}{})
}

func runGenGo(cmd *cobra.Command, args []string) error {
	opts := make(map[string]interface{})
	if genGoPackage != "" {
		opts["package"] = genGoPackage
	}
	return runGenerator("go", opts)
}

func runGenElixir(cmd *cobra.Command, args []string) error {
	opts := make(map[string]interface{})
	if genElixirModule != "" {
		opts["module"] = genElixirModule
	}
	return runGenerator("elixir", opts)
}

// runGenerator is a generic function to run any generator
func runGenerator(name string, opts map[string]interface{}) error {
	// Load config
	cfg, err := config.Load(GetConfigFile())
	if err != nil {
		return err
	}

	// Get generator config
	genCfg, ok := cfg.Generate[name]
	if !ok {
		genCfg = config.GenConfig{
			Enabled: true,
			Output:  fmt.Sprintf("generated/%s", getDefaultOutput(name)),
			Options: make(map[string]interface{}),
		}
	}

	// Override output if specified
	if genOutput != "" {
		genCfg.Output = genOutput
	}

	// Merge options
	if genCfg.Options == nil {
		genCfg.Options = make(map[string]interface{})
	}
	for k, v := range opts {
		genCfg.Options[k] = v
	}

	// Validate output path
	if genCfg.Output == "" {
		return fmt.Errorf("no output file specified (use --output or configure in platosl.yaml)")
	}

	PrintVerbose("Generating %s to: %s", name, genCfg.Output)

	// Load and validate schemas
	val, err := loadAndValidateSchemas(cfg, name)
	if err != nil {
		return err
	}

	// Get generator
	gen, err := generator.Get(name)
	if err != nil {
		e := errors.Wrap(errors.ErrorTypeInternal, err, fmt.Sprintf("%s generator not registered", name))
		e = e.WithSuggestion("This is an internal error. Please report this issue")
		PrintError(e.Format())
		return e
	}

	// Create generator context
	ctx := generator.NewContext(val, cfg, genCfg)

	// Validate
	PrintVerbose("Validating generator requirements")
	if err := gen.Validate(ctx); err != nil {
		e := errors.Wrap(errors.ErrorTypeValidation, err, fmt.Sprintf("%s generator validation failed", name))
		e = e.WithSuggestion(fmt.Sprintf("The schema structure may not be compatible with %s generation", name))
		PrintError(e.Format())
		return e
	}

	// Generate
	PrintVerbose("Generating %s code", name)
	output, err := gen.Generate(ctx)
	if err != nil {
		e := errors.Wrap(errors.ErrorTypeGeneration, err, fmt.Sprintf("%s generation failed", name))
		e = e.WithSuggestion("Check that your schema definitions are valid and exportable")
		PrintError(e.Format())
		return e
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(genCfg.Output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		e := errors.Wrap(errors.ErrorTypeFileSystem, err, fmt.Sprintf("failed to create output directory: %s", outputDir))
		e = e.WithSuggestion("Check that you have write permissions for the output directory")
		PrintError(e.Format())
		return e
	}

	// Write output
	if err := os.WriteFile(genCfg.Output, output, 0644); err != nil {
		e := errors.Wrap(errors.ErrorTypeFileSystem, err, fmt.Sprintf("failed to write output file: %s", genCfg.Output))
		e = e.WithSuggestion("Check that you have write permissions for the output file")
		PrintError(e.Format())
		return e
	}

	// Success
	stats := fmt.Sprintf("%d bytes", len(output))
	PrintSuccess("Generated %s: %s (%s)", name, filepath.Base(genCfg.Output), stats)

	return nil
}

func getDefaultOutput(generatorName string) string {
	switch generatorName {
	case "typescript":
		return "types.ts"
	case "zod":
		return "schemas.ts"
	case "jsonschema":
		return "schema.json"
	case "go":
		return "types.go"
	case "elixir":
		return "types.ex"
	default:
		return "output.txt"
	}
}

// validateSchemas performs validation on loaded schemas and returns structured errors
func validateSchemas(val cue.Value, generatorName string) []*errors.Error {
	var errs []*errors.Error

	// Create validator
	validator := platoCue.NewValidator(false)
	result := validator.Validate(val)

	if !result.Valid {
		for _, valErr := range result.Errors {
			err := errors.New(errors.ErrorTypeValidation, valErr.Message).
				WithLocation(valErr.File, valErr.Line, valErr.Column)

			if valErr.Suggestion != "" {
				err = err.WithSuggestion(valErr.Suggestion)
			} else if valErr.Path != "" {
				err = err.WithSuggestion(fmt.Sprintf("Check field '%s' in your schema", valErr.Path))
			}

			errs = append(errs, err)
		}
	}

	return errs
}

// loadAndValidateSchemas loads schemas and performs validation
func loadAndValidateSchemas(cfg *config.Config, generatorName string) (cue.Value, error) {
	loader := platoCue.NewLoader()

	var allPaths []string
	for _, schemaPath := range cfg.Schemas {
		absPath, err := filepath.Abs(schemaPath)
		if err != nil {
			e := errors.Wrap(errors.ErrorTypeFileSystem, err, fmt.Sprintf("failed to resolve schema path: %s", schemaPath))
			e = e.WithSuggestion("Verify that the schema path in platosl.yaml exists and is accessible")
			PrintError(e.Format())
			return cue.Value{}, e
		}

		// Check if path exists
		if _, err := os.Stat(absPath); err != nil {
			if os.IsNotExist(err) {
				e := errors.New(errors.ErrorTypeFileSystem, fmt.Sprintf("schema path not found: %s", schemaPath))
				e = e.WithSuggestion("Create the directory or update the 'schemas' section in platosl.yaml")
				PrintError(e.Format())
				return cue.Value{}, e
			}
		}

		allPaths = append(allPaths, absPath)
	}

	if len(allPaths) == 0 {
		e := errors.New(errors.ErrorTypeConfig, "no schema paths configured in platosl.yaml")
		e = e.WithSuggestion("Add schema directories to the 'schemas' section in platosl.yaml")
		PrintError(e.Format())
		return cue.Value{}, e
	}

	PrintVerbose("Loading %d schema path(s) for %s generation", len(allPaths), generatorName)

	// Load all schemas
	val, err := loader.LoadPaths(allPaths)
	if err != nil {
		// Provide context-specific suggestions
		suggestion := "Check your CUE files for syntax errors. Run 'cue vet' directly for more details"
		if strings.Contains(err.Error(), "cannot use absolute directory") {
			suggestion = "CUE module configuration issue. Try using relative paths in platosl.yaml or ensure you have a cue.mod directory"
		} else if strings.Contains(err.Error(), "import failed") {
			suggestion = "Check that all imported packages are available in your cue.mod directory"
		} else if strings.Contains(err.Error(), "cannot find package") {
			suggestion = "Verify that the schema paths in platosl.yaml point to valid CUE packages"
		}

		e := errors.Wrap(errors.ErrorTypeValidation, err, "failed to load CUE schemas")
		e = e.WithSuggestion(suggestion)
		PrintError(e.Format())
		return cue.Value{}, e
	}

	// Validate schemas
	validationErrors := validateSchemas(val, generatorName)
	if len(validationErrors) > 0 {
		PrintError("Schema validation failed with %d error(s):\n", len(validationErrors))
		for _, err := range validationErrors {
			PrintError(err.Format())
			fmt.Fprintln(os.Stderr)
		}
		return cue.Value{}, fmt.Errorf("schema validation failed")
	}

	PrintVerbose("✓ All schemas validated successfully")

	return val, nil
}

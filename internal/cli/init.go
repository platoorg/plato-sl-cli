package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/platoorg/plato-sl-cli/internal/config"
)

var (
	initBase       string
	initName       string
	initGenerators string
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new PlatoSL project or update generator configuration",
	Long: `Initialize a new PlatoSL project with a platosl.yaml configuration file
and directory structure. If a platosl.yaml already exists, this command will
allow you to update the generator configuration.

If a directory is specified, the project will be initialized there.
Otherwise, it will be initialized in the current directory.

By default, the command will interactively prompt you to select which
generators to enable. You can also specify generators using the --generators flag.

Available generators:
  - typescript  : TypeScript interfaces
  - zod         : Zod validation schemas
  - go          : Go structs
  - jsonschema  : JSON Schema
  - elixir      : Elixir typespecs

Examples:
  # Initialize interactively (will prompt for generator selection)
  platosl init

  # Re-run in an existing project to add/remove generators
  platosl init

  # Initialize with specific generators (non-interactive)
  platosl init --generators typescript,go,jsonschema

  # Update generators in existing project (non-interactive)
  platosl init --generators typescript,zod,jsonschema,go,elixir`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initBase, "base", "", "base schema to import (e.g., platosl.org/base/address/us@v1)")
	initCmd.Flags().StringVar(&initName, "name", "", "project name (defaults to directory name)")
	initCmd.Flags().StringVar(&initGenerators, "generators", "typescript,zod", "comma-separated list of generators to enable (typescript,zod,jsonschema,go,elixir)")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Get absolute path
	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve directory path: %w", err)
	}

	// Check if directory exists, create if not
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		PrintVerbose("Creating directory: %s", absDir)
		if err := os.MkdirAll(absDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Check if already initialized
	configPath := filepath.Join(absDir, "platosl.yaml")
	existingConfig := config.Exists(configPath)

	var cfg *config.Config
	var projectName string
	var currentGenerators []string

	if existingConfig {
		// Load existing config
		PrintInfo("Found existing platosl.yaml - re-initializing project")
		var err error
		cfg, err = config.Load(configPath)
		if err != nil {
			return fmt.Errorf("failed to load existing config: %w", err)
		}

		projectName = cfg.Name

		// Get currently enabled generators
		for genName, genCfg := range cfg.Generate {
			if genCfg.Enabled {
				currentGenerators = append(currentGenerators, genName)
			}
		}

		PrintVerbose("Current generators: %s", strings.Join(currentGenerators, ", "))
	} else {
		// Determine project name for new project
		projectName = initName
		if projectName == "" {
			projectName = filepath.Base(absDir)
		}
		PrintVerbose("Initializing new project: %s", projectName)
	}

	// Parse selected generators
	var selectedGenerators []string

	// Check if generators flag was explicitly provided
	generatorsFlag := cmd.Flags().Lookup("generators")
	flagWasSet := generatorsFlag != nil && generatorsFlag.Changed

	if flagWasSet {
		// Use the provided generators from flag
		selectedGenerators = strings.Split(initGenerators, ",")
		// Trim whitespace
		for i, gen := range selectedGenerators {
			selectedGenerators[i] = strings.TrimSpace(gen)
		}
		PrintVerbose("Enabling generators: %s", strings.Join(selectedGenerators, ", "))
	} else {
		// Interactive mode - prompt user to select generators
		availableGenerators := []string{"typescript", "zod", "go", "jsonschema", "elixir"}
		defaultGenerators := currentGenerators
		if len(defaultGenerators) == 0 {
			defaultGenerators = []string{"typescript", "zod"}
		}

		message := "Select generators to enable:"
		if existingConfig {
			message = "Update generator selection (currently: " + strings.Join(currentGenerators, ", ") + "):"
		}

		prompt := &survey.MultiSelect{
			Message: message,
			Options: availableGenerators,
			Default: defaultGenerators,
			Help:    "Use space to select/deselect, enter to confirm. Multiple generators can be selected.",
		}

		if err := survey.AskOne(prompt, &selectedGenerators, survey.WithValidator(survey.Required)); err != nil {
			return fmt.Errorf("generator selection cancelled or failed: %w", err)
		}

		PrintInfo("Selected generators: %s", strings.Join(selectedGenerators, ", "))
	}

	// Create or update config with selected generators
	if existingConfig {
		// Update existing config with new generator selection
		cfg = config.UpdateGenerators(cfg, selectedGenerators)
	} else {
		// Create new config with selected generators
		cfg = config.DefaultWithGenerators(projectName, selectedGenerators)
	}

	// Add base schema if specified
	if initBase != "" {
		PrintVerbose("Adding base schema: %s", initBase)
		cfg.Imports = append(cfg.Imports, initBase)
	}

	// Create directory structure
	dirs := []string{
		filepath.Join(absDir, "schemas"),
		filepath.Join(absDir, "generated"),
	}

	for _, dir := range dirs {
		PrintVerbose("Creating directory: %s", filepath.Base(dir))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Save config
	PrintVerbose("Writing platosl.yaml")
	if err := config.Save(configPath, cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Create example schema
	exampleSchema := filepath.Join(absDir, "schemas", "example.cue")
	exampleContent := `package schemas

// Example schema
#Person: {
	name!: string
	email!: string & =~"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	age?: int & >=0 & <=150
}
`
	PrintVerbose("Creating example schema: schemas/example.cue")
	if err := os.WriteFile(exampleSchema, []byte(exampleContent), 0644); err != nil {
		return fmt.Errorf("failed to create example schema: %w", err)
	}

	// Success message
	if existingConfig {
		PrintSuccess("Updated PlatoSL project configuration: %s", projectName)
		PrintInfo("")
		PrintInfo("Generators enabled: %s", strings.Join(selectedGenerators, ", "))
		PrintInfo("")
		PrintInfo("Next steps:")
		PrintInfo("  1. Run 'platosl validate' to validate schemas")
		PrintInfo("  2. Run 'platosl build' to generate code for all enabled generators")
	} else {
		PrintSuccess("Initialized PlatoSL project: %s", projectName)
		PrintInfo("")
		PrintInfo("Created:")
		PrintInfo("  platosl.yaml        - Configuration file")
		PrintInfo("  schemas/            - Schema directory")
		PrintInfo("  schemas/example.cue - Example schema")
		PrintInfo("  generated/          - Generated code output")
		PrintInfo("")
		PrintInfo("Next steps:")
		PrintInfo("  1. Edit schemas/example.cue or add your own schemas")
		PrintInfo("  2. Run 'platosl validate' to validate schemas")
		PrintInfo("  3. Run 'platosl gen typescript' to generate TypeScript types")
	}

	return nil
}

package cli

import (
	"github.com/spf13/cobra"
	"github.com/platoorg/platosl-cli/internal/config"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Validate schemas and generate all enabled targets",
	Long: `Build validates all CUE schemas and generates code for all enabled targets
configured in platosl.yaml.

This is equivalent to running 'platosl validate' followed by generating all
enabled generators.`,
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load(GetConfigFile())
	if err != nil {
		return err
	}

	PrintInfo("Building project: %s", cfg.Name)
	PrintInfo("")

	// Step 1: Validate
	PrintInfo("Step 1: Validating schemas...")
	if err := runValidate(cmd, []string{}); err != nil {
		return err
	}
	PrintInfo("")

	// Step 2: Generate all
	PrintInfo("Step 2: Generating code...")
	if err := runGenAll(cfg); err != nil {
		return err
	}

	PrintInfo("")
	PrintSuccess("Build complete")
	return nil
}

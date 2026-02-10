package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	platoCue "github.com/platoorg/plato-sl-cli/internal/cue"
)

var (
	infoFormat string
)

var infoCmd = &cobra.Command{
	Use:   "info <schema>",
	Short: "Show schema information",
	Long: `Show detailed information about a CUE schema including fields, types,
and definitions.`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVar(&infoFormat, "format", "text", "output format (text, json, yaml)")
}

func runInfo(cmd *cobra.Command, args []string) error {
	schemaPath := args[0]

	// Resolve path
	absPath, err := filepath.Abs(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	PrintVerbose("Loading schema: %s", schemaPath)

	// Load schema
	loader := platoCue.NewLoader()
	val, err := loader.LoadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Introspect schema
	info, err := platoCue.Introspect(val)
	if err != nil {
		return fmt.Errorf("failed to introspect schema: %w", err)
	}

	// Format output
	switch infoFormat {
	case "json":
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format as JSON: %w", err)
		}
		fmt.Println(string(data))

	case "yaml":
		data, err := yaml.Marshal(info)
		if err != nil {
			return fmt.Errorf("failed to format as YAML: %w", err)
		}
		fmt.Print(string(data))

	case "text":
		fallthrough
	default:
		fmt.Print(platoCue.FormatSchemaInfo(info))
	}

	return nil
}

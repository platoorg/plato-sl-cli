package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "platosl",
	Short: "PlatoSL - Schema language for content validation",
	Long: `PlatoSL is a CLI tool for managing CUE-based schemas for content validation.
It provides commands for initialization, validation, and code generation from
CUE schemas to TypeScript, JSON Schema, Go, and Elixir.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is platosl.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

// GetConfigFile returns the config file path
func GetConfigFile() string {
	if cfgFile != "" {
		return cfgFile
	}
	return "platosl.yaml"
}

// PrintError prints an error message with formatting
func PrintError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "✗ "+msg+"\n", args...)
}

// PrintSuccess prints a success message with formatting
func PrintSuccess(msg string, args ...interface{}) {
	fmt.Printf("✓ "+msg+"\n", args...)
}

// PrintInfo prints an info message
func PrintInfo(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// PrintVerbose prints a message only in verbose mode
func PrintVerbose(msg string, args ...interface{}) {
	if verbose {
		fmt.Printf("  "+msg+"\n", args...)
	}
}

package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the CLI
	// This can be set at build time using: -ldflags "-X github.com/platoorg/plato-sl-cli/internal/cli.Version=v1.0.0"
	Version = "plato-sl-cli-0.0.2"

	// Commit is the git commit hash
	// This can be set at build time using: -ldflags "-X github.com/platoorg/plato-sl-cli/internal/cli.Commit=abc123"
	Commit = "unknown"

	// BuildDate is the date the binary was built
	// This can be set at build time using: -ldflags "-X github.com/platoorg/plato-sl-cli/internal/cli.BuildDate=2024-01-01"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display version information for the PlatoSL CLI including version, commit hash, build date, and Go version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Commit:     %s\n", Commit)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

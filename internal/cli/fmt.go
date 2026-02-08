package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"platosl.org/cmd/platosl/internal/config"
)

var (
	fmtCheck bool
	fmtWrite bool
)

var fmtCmd = &cobra.Command{
	Use:   "fmt [file or directory]",
	Short: "Format CUE files",
	Long: `Format CUE files using 'cue fmt'.

If a file or directory is specified, formats only that path.
Otherwise, formats all schema paths from platosl.yaml.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFmt,
}

func init() {
	rootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().BoolVar(&fmtCheck, "check", false, "check if files are formatted (exit 1 if not)")
	fmtCmd.Flags().BoolVarP(&fmtWrite, "write", "w", true, "write result to (source) file")
}

func runFmt(cmd *cobra.Command, args []string) error {
	// Determine what to format
	var paths []string

	if len(args) > 0 {
		// Format specific path
		path := args[0]
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}

		paths = []string{absPath}
	} else {
		// Load config and format configured paths
		cfg, err := config.Load(GetConfigFile())
		if err != nil {
			return err
		}

		// Collect all schema paths
		for _, schemaPath := range cfg.Schemas {
			absPath, err := filepath.Abs(schemaPath)
			if err != nil {
				PrintError("Failed to resolve path %s: %v", schemaPath, err)
				continue
			}
			paths = append(paths, absPath)
		}

		if len(paths) == 0 {
			return fmt.Errorf("no schema paths configured in platosl.yaml")
		}
	}

	// Check if 'cue' command is available
	if _, err := exec.LookPath("cue"); err != nil {
		return fmt.Errorf("'cue' command not found\n\nInstall CUE: https://cuelang.org/docs/install/")
	}

	// Format each path
	formatted := 0
	for _, path := range paths {
		PrintVerbose("Formatting: %s", path)

		// Build cue fmt command
		cmdArgs := []string{"fmt"}
		if fmtCheck {
			// Use diff mode to check
			cmdArgs = append(cmdArgs, "-d")
		}
		cmdArgs = append(cmdArgs, path)

		// Run cue fmt
		cueCmd := exec.Command("cue", cmdArgs...)
		output, err := cueCmd.CombinedOutput()

		if err != nil {
			if fmtCheck {
				// Check mode - show diff
				fmt.Print(string(output))
				return fmt.Errorf("files not formatted")
			}
			return fmt.Errorf("failed to format %s: %w\n%s", path, err, string(output))
		}

		if fmtCheck && len(output) > 0 {
			// Has diff output
			fmt.Print(string(output))
			return fmt.Errorf("files not formatted")
		}

		formatted++
	}

	if fmtCheck {
		PrintSuccess("All files formatted correctly")
	} else {
		PrintSuccess("Formatted %d path(s)", formatted)
	}

	return nil
}

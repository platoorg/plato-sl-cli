package cue

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// Loader handles loading CUE files and directories
type Loader struct {
	ctx *cue.Context
}

// NewLoader creates a new CUE loader
func NewLoader() *Loader {
	return &Loader{
		ctx: cuecontext.New(),
	}
}

// LoadFile loads a single CUE file
func (l *Loader) LoadFile(path string) (cue.Value, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return cue.Value{}, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	val := l.ctx.CompileBytes(data, cue.Filename(path))
	if err := val.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("failed to compile %s: %w", path, err)
	}

	return val, nil
}

// LoadDir loads all CUE files from a directory
func (l *Loader) LoadDir(dir string) (cue.Value, error) {
	// Check if directory exists
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return cue.Value{}, fmt.Errorf("directory not found: %s", dir)
		}
		return cue.Value{}, fmt.Errorf("failed to stat directory: %w", err)
	}
	if !info.IsDir() {
		return cue.Value{}, fmt.Errorf("not a directory: %s", dir)
	}

	// Try module-based loading first
	moduleRoot := findModuleRoot(dir)
	hasModule := moduleRoot != "" && dirExists(filepath.Join(moduleRoot, "cue.mod"))

	if hasModule {
		// Use load.Instances for module-based loading
		loadPath := dir
		if !filepath.IsAbs(dir) && !strings.HasPrefix(dir, "./") && !strings.HasPrefix(dir, "../") {
			loadPath = "./" + dir
		}

		cfg := &load.Config{
			ModuleRoot: moduleRoot,
		}
		buildInstances := load.Instances([]string{loadPath}, cfg)
		if len(buildInstances) > 0 && buildInstances[0].Err == nil {
			inst := buildInstances[0]
			val := l.ctx.BuildInstance(inst)
			if err := val.Err(); err == nil {
				return val, nil
			}
		}
	}

	// Fallback: Load individual CUE files from directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return cue.Value{}, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var values []cue.Value
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".cue") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return cue.Value{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		val := l.ctx.CompileBytes(data, cue.Filename(filePath))
		if err := val.Err(); err != nil {
			return cue.Value{}, fmt.Errorf("failed to compile %s: %w", filePath, err)
		}

		values = append(values, val)
	}

	if len(values) == 0 {
		return cue.Value{}, fmt.Errorf("no CUE files found in %s", dir)
	}

	// Unify all values from the directory
	result := values[0]
	for i := 1; i < len(values); i++ {
		result = result.Unify(values[i])
		if err := result.Err(); err != nil {
			return cue.Value{}, fmt.Errorf("failed to unify CUE files in %s: %w", dir, err)
		}
	}

	return result, nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// LoadPaths loads CUE files from multiple paths (files or directories)
func (l *Loader) LoadPaths(paths []string) (cue.Value, error) {
	if len(paths) == 0 {
		return cue.Value{}, fmt.Errorf("no paths provided")
	}

	var values []cue.Value

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return cue.Value{}, fmt.Errorf("failed to stat %s: %w", path, err)
		}

		var val cue.Value
		if info.IsDir() {
			val, err = l.LoadDir(path)
		} else {
			val, err = l.LoadFile(path)
		}

		if err != nil {
			return cue.Value{}, err
		}

		values = append(values, val)
	}

	// Unify all values
	result := values[0]
	for i := 1; i < len(values); i++ {
		result = result.Unify(values[i])
		if err := result.Err(); err != nil {
			return cue.Value{}, fmt.Errorf("failed to unify values: %w", err)
		}
	}

	return result, nil
}

// Context returns the CUE context
func (l *Loader) Context() *cue.Context {
	return l.ctx
}

// ExpandGlob expands glob patterns to file paths
func ExpandGlob(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern %s: %w", pattern, err)
	}
	return matches, nil
}

// findModuleRoot searches for cue.mod directory starting from the given directory
func findModuleRoot(dir string) string {
	// Try to find cue.mod in current directory or parents
	current := dir
	for {
		modPath := filepath.Join(current, "cue.mod")
		if _, err := os.Stat(modPath); err == nil {
			return current
		}

		// Go up one directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached root
			break
		}
		current = parent
	}

	// Return current working directory as fallback
	cwd, _ := os.Getwd()
	return cwd
}

package config

// Config represents the platosl.yaml configuration
type Config struct {
	Version    string              `yaml:"version"`
	Name       string              `yaml:"name"`
	Imports    []string            `yaml:"imports,omitempty"`
	Schemas    []string            `yaml:"schemas"`
	Validation ValidationConfig    `yaml:"validation"`
	Generate   map[string]GenConfig `yaml:"generate"`
}

// ValidationConfig holds validation options
type ValidationConfig struct {
	Strict        bool `yaml:"strict"`
	FailOnWarning bool `yaml:"failOnWarning"`
}

// GenConfig holds generator-specific configuration
type GenConfig struct {
	Enabled bool                   `yaml:"enabled"`
	Output  string                 `yaml:"output"`
	Options map[string]interface{} `yaml:"options,omitempty"`
}

// TypeScriptOptions holds TypeScript-specific options
type TypeScriptOptions struct {
	Zod bool `yaml:"zod"`
}

// GoOptions holds Go-specific options
type GoOptions struct {
	Package string `yaml:"package"`
}

// ElixirOptions holds Elixir-specific options
type ElixirOptions struct {
	Module string `yaml:"module"`
}

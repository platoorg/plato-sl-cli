package config

// Default returns a default configuration with all generators enabled
func Default(name string) *Config {
	return DefaultWithGenerators(name, []string{"typescript", "zod", "jsonschema", "go", "elixir"})
}

// DefaultWithGenerators returns a configuration with specific generators enabled
func DefaultWithGenerators(name string, generators []string) *Config {
	if name == "" {
		name = "my-project"
	}

	cfg := &Config{
		Version: "v1",
		Name:    name,
		Imports: []string{},
		Schemas: []string{"schemas/"},
		Validation: ValidationConfig{
			Strict:        true,
			FailOnWarning: false,
		},
		Generate: make(map[string]GenConfig),
	}

	// Add requested generators
	for _, gen := range generators {
		switch gen {
		case "typescript":
			cfg.Generate["typescript"] = GenConfig{
				Enabled: true,
				Output:  "generated/types.ts",
			}
		case "zod":
			cfg.Generate["zod"] = GenConfig{
				Enabled: true,
				Output:  "generated/schemas.ts",
			}
		case "jsonschema":
			cfg.Generate["jsonschema"] = GenConfig{
				Enabled: true,
				Output:  "generated/schema.json",
			}
		case "go":
			cfg.Generate["go"] = GenConfig{
				Enabled: true,
				Output:  "generated/types.go",
				Options: map[string]interface{}{
					"package": "types",
				},
			}
		case "elixir":
			cfg.Generate["elixir"] = GenConfig{
				Enabled: true,
				Output:  "generated/types.ex",
				Options: map[string]interface{}{
					"module": "MyApp.Types",
				},
			}
		}
	}

	return cfg
}

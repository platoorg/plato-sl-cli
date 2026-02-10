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

// UpdateGenerators updates an existing configuration with new generator selection
func UpdateGenerators(cfg *Config, generators []string) *Config {
	// Create a map of selected generators for quick lookup
	selectedMap := make(map[string]bool)
	for _, gen := range generators {
		selectedMap[gen] = true
	}

	// Update existing generators and disable those not selected
	allGenerators := []string{"typescript", "zod", "jsonschema", "go", "elixir"}
	for _, gen := range allGenerators {
		if existingCfg, exists := cfg.Generate[gen]; exists {
			// Generator exists in config - update enabled status
			existingCfg.Enabled = selectedMap[gen]
			cfg.Generate[gen] = existingCfg
		} else if selectedMap[gen] {
			// Generator doesn't exist but is selected - add it with defaults
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
	}

	return cfg
}

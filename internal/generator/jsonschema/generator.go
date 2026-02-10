package jsonschema

import (
	"encoding/json"
	"fmt"

	"github.com/platoorg/plato-sl-cli/internal/generator"
)

// Generator generates JSON Schema from CUE
type Generator struct{}

// NewGenerator creates a new JSON Schema generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Name returns the generator name
func (g *Generator) Name() string {
	return "jsonschema"
}

// Generate generates JSON Schema
func (g *Generator) Generate(ctx *generator.Context) ([]byte, error) {
	// Use CUE's built-in JSON marshaling
	data, err := ctx.Value.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Parse and wrap in JSON Schema format
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create JSON Schema wrapper
	schema := map[string]interface{}{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         fmt.Sprintf("https://platosl.org/schemas/%s", ctx.Config.Name),
		"title":       ctx.Config.Name,
		"type":        "object",
		"properties":  obj,
		"definitions": extractDefinitions(obj),
	}

	// Pretty-print JSON
	output, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %w", err)
	}

	return output, nil
}

// Validate validates the generator context
func (g *Generator) Validate(ctx *generator.Context) error {
	if err := ctx.Value.Err(); err != nil {
		return fmt.Errorf("invalid CUE value: %w", err)
	}
	return nil
}

// extractDefinitions extracts definitions from the object
func extractDefinitions(obj map[string]interface{}) map[string]interface{} {
	defs := make(map[string]interface{})

	for key, val := range obj {
		// CUE definitions start with #
		if len(key) > 0 && key[0] == '#' {
			defs[key[1:]] = val
			delete(obj, key)
		}
	}

	return defs
}

func init() {
	// Register the generator
	generator.Register(NewGenerator())
}

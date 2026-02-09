package generator

import (
	"cuelang.org/go/cue"
	"github.com/platoorg/platosl-cli/internal/config"
)

// Generator is the interface that all code generators must implement
type Generator interface {
	// Name returns the generator's name (e.g., "typescript", "go")
	Name() string

	// Generate generates code from a CUE value
	Generate(ctx *Context) ([]byte, error)

	// Validate validates that the generator can process the given context
	Validate(ctx *Context) error
}

// Context holds the context for code generation
type Context struct {
	// Value is the CUE value to generate code from
	Value cue.Value

	// Config is the project configuration
	Config *config.Config

	// GeneratorConfig is the generator-specific configuration
	GeneratorConfig config.GenConfig

	// Options contains additional generator options
	Options map[string]interface{}
}

// NewContext creates a new generator context
func NewContext(value cue.Value, cfg *config.Config, genCfg config.GenConfig) *Context {
	ctx := &Context{
		Value:           value,
		Config:          cfg,
		GeneratorConfig: genCfg,
		Options:         make(map[string]interface{}),
	}

	// Merge generator options
	if genCfg.Options != nil {
		for k, v := range genCfg.Options {
			ctx.Options[k] = v
		}
	}

	return ctx
}

// GetOption retrieves an option value
func (c *Context) GetOption(key string) (interface{}, bool) {
	val, ok := c.Options[key]
	return val, ok
}

// GetStringOption retrieves a string option
func (c *Context) GetStringOption(key string, defaultVal string) string {
	if val, ok := c.Options[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultVal
}

// GetBoolOption retrieves a boolean option
func (c *Context) GetBoolOption(key string, defaultVal bool) bool {
	if val, ok := c.Options[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultVal
}

// GetIntOption retrieves an integer option
func (c *Context) GetIntOption(key string, defaultVal int) int {
	if val, ok := c.Options[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
		// Try float64 (JSON numbers)
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return defaultVal
}

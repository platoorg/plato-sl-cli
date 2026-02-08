package generator

import (
	"fmt"
	"sync"
)

// Registry manages available generators
type Registry struct {
	mu         sync.RWMutex
	generators map[string]Generator
}

var (
	// DefaultRegistry is the global generator registry
	DefaultRegistry = NewRegistry()
)

// NewRegistry creates a new generator registry
func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]Generator),
	}
}

// Register registers a generator
func (r *Registry) Register(gen Generator) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := gen.Name()
	if _, exists := r.generators[name]; exists {
		return fmt.Errorf("generator %s already registered", name)
	}

	r.generators[name] = gen
	return nil
}

// Get retrieves a generator by name
func (r *Registry) Get(name string) (Generator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	gen, exists := r.generators[name]
	if !exists {
		return nil, fmt.Errorf("generator %s not found", name)
	}

	return gen, nil
}

// List returns all registered generator names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.generators))
	for name := range r.generators {
		names = append(names, name)
	}
	return names
}

// Register is a convenience function that registers a generator in the default registry
func Register(gen Generator) error {
	return DefaultRegistry.Register(gen)
}

// Get is a convenience function that retrieves a generator from the default registry
func Get(name string) (Generator, error) {
	return DefaultRegistry.Get(name)
}

// List is a convenience function that lists all generators in the default registry
func List() []string {
	return DefaultRegistry.List()
}

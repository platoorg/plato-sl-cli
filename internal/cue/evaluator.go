package cue

import (
	"fmt"

	"cuelang.org/go/cue"
)

// Evaluator evaluates CUE expressions
type Evaluator struct {
	loader *Loader
}

// NewEvaluator creates a new evaluator
func NewEvaluator(loader *Loader) *Evaluator {
	return &Evaluator{loader: loader}
}

// Evaluate evaluates a CUE value to its concrete form
func (e *Evaluator) Evaluate(val cue.Value) (interface{}, error) {
	// Decode the value to a Go interface
	var result interface{}
	if err := val.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode value: %w", err)
	}
	return result, nil
}

// EvaluateJSON evaluates a CUE value and returns JSON bytes
func (e *Evaluator) EvaluateJSON(val cue.Value) ([]byte, error) {
	// Use CUE's built-in JSON marshaling
	return val.MarshalJSON()
}

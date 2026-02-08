package cue

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
)

// SchemaInfo holds information about a CUE schema
type SchemaInfo struct {
	Fields      []FieldInfo
	Definitions []string
}

// FieldInfo holds information about a field
type FieldInfo struct {
	Name     string
	Type     string
	Optional bool
	Path     string
}

// Introspect extracts schema information from a CUE value
func Introspect(val cue.Value) (*SchemaInfo, error) {
	info := &SchemaInfo{
		Fields:      []FieldInfo{},
		Definitions: []string{},
	}

	// Walk the value structure
	iter, err := val.Fields(cue.Optional(true), cue.Definitions(true))
	if err != nil {
		return nil, fmt.Errorf("failed to iterate fields: %w", err)
	}

	for iter.Next() {
		label := iter.Selector().String()
		value := iter.Value()

		// Check if it's a definition
		if strings.HasPrefix(label, "#") {
			info.Definitions = append(info.Definitions, label)
		}

		// Extract field info
		fieldInfo := FieldInfo{
			Name:     label,
			Type:     inferType(value),
			Optional: iter.IsOptional(),
			Path:     iter.Selector().String(),
		}

		info.Fields = append(info.Fields, fieldInfo)
	}

	return info, nil
}

// inferType infers the CUE type as a string
func inferType(val cue.Value) string {
	kind := val.IncompleteKind()

	switch {
	case kind&cue.StringKind != 0:
		return "string"
	case kind&cue.IntKind != 0:
		return "int"
	case kind&cue.FloatKind != 0:
		return "float"
	case kind&cue.NumberKind != 0:
		return "number"
	case kind&cue.BoolKind != 0:
		return "bool"
	case kind&cue.ListKind != 0:
		return "list"
	case kind&cue.StructKind != 0:
		return "struct"
	default:
		return "unknown"
	}
}

// FormatSchemaInfo formats schema info as a string
func FormatSchemaInfo(info *SchemaInfo) string {
	var b strings.Builder

	if len(info.Definitions) > 0 {
		b.WriteString("Definitions:\n")
		for _, def := range info.Definitions {
			fmt.Fprintf(&b, "  %s\n", def)
		}
		b.WriteString("\n")
	}

	if len(info.Fields) > 0 {
		b.WriteString("Fields:\n")
		for _, field := range info.Fields {
			optional := ""
			if field.Optional {
				optional = " (optional)"
			}
			fmt.Fprintf(&b, "  %s: %s%s\n", field.Name, field.Type, optional)
		}
	}

	return b.String()
}

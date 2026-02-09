# PlatoSL CLI Documentation

The PlatoSL CLI provides commands for managing CUE schemas and generating code in multiple languages.

## Installation

```bash
# Install from source
go install github.com/platoorg/platosl-cli@latest

# Or build locally
git clone https://github.com/platoOrg/platoSl.git
cd platoSl
make install
```

## Quick Start

```bash
# Initialize a new project
platosl init my-project --base platosl.org/base/address/us

# Validate schemas
platosl validate

# Generate TypeScript types with Zod validation
platosl gen typescript --zod

# Generate all configured targets
platosl build
```

## Commands

### `platosl init`

Initialize a new PlatoSL project with configuration and directory structure.

```bash
platosl init [directory] [flags]

Flags:
  --base string    Base schema to import (e.g., platosl.org/base/address/us@v1)
  --name string    Project name (defaults to directory name)
```

**Example:**
```bash
platosl init my-project --base platosl.org/base/address/us
cd my-project
```

This creates:
- `platosl.yaml` - Configuration file
- `schemas/` - Schema directory
- `schemas/example.cue` - Example schema
- `generated/` - Generated code output directory

---

### `platosl validate`

Validate CUE schemas for correctness and completeness.

```bash
platosl validate [file or directory] [flags]

Flags:
  --strict    Strict validation (requires all fields to be concrete)
```

**Examples:**
```bash
# Validate all schemas from config
platosl validate

# Validate specific file
platosl validate schemas/address.cue

# Validate directory
platosl validate schemas/

# Strict mode
platosl validate --strict
```

---

### `platosl gen`

Generate code from CUE schemas to various target languages.

#### `platosl gen typescript`

Generate TypeScript interfaces and Zod validation schemas.

```bash
platosl gen typescript [flags]

Flags:
  -o, --output string    Output file path
      --zod              Generate Zod schemas for runtime validation
```

**Example:**
```bash
# Generate TypeScript with Zod
platosl gen typescript --zod --output src/types.ts

# Using config defaults
platosl gen typescript
```

**Output:**
```typescript
// Generated types
export interface Person {
  name: string;
  email: string;
  age?: number;
}

// Zod schemas (if --zod is enabled)
export const PersonSchema = z.object({
  name: z.string(),
  email: z.string(),
  age: z.number().int().optional(),
});
```

#### `platosl gen jsonschema`

Generate JSON Schema (draft 2020-12).

```bash
platosl gen jsonschema [flags]

Flags:
  -o, --output string    Output file path
```

**Example:**
```bash
platosl gen jsonschema --output schema.json
```

#### `platosl gen go`

Generate Go structs with JSON tags.

```bash
platosl gen go [flags]

Flags:
  -o, --output string       Output file path
      --package string      Go package name
```

**Example:**
```bash
platosl gen go --output pkg/types/types.go --package types
```

**Output:**
```go
package types

type Person struct {
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Age   *int    `json:"age,omitempty"`
}
```

#### `platosl gen elixir`

Generate Elixir typespecs and structs.

```bash
platosl gen elixir [flags]

Flags:
  -o, --output string     Output file path
      --module string     Elixir module name
```

**Example:**
```bash
platosl gen elixir --output lib/types.ex --module MyApp.Types
```

**Output:**
```elixir
defmodule MyApp.Types do
  @type person() :: %__MODULE__.Person{
    name: String.t(),
    email: String.t(),
    age: integer() | nil
  }

  defstruct [:name, :email, :age]
end
```

---

### `platosl build`

Validate schemas and generate all enabled targets.

```bash
platosl build
```

This command:
1. Validates all schemas
2. Generates code for all enabled generators in `platosl.yaml`

Equivalent to running `platosl validate` followed by generating all targets.

---

### `platosl fmt`

Format CUE files using `cue fmt`.

```bash
platosl fmt [file or directory] [flags]

Flags:
      --check       Check if files are formatted (exit 1 if not)
  -w, --write       Write result to source file (default true)
```

**Examples:**
```bash
# Format all schemas
platosl fmt

# Format specific file
platosl fmt schemas/address.cue

# Check formatting in CI
platosl fmt --check
```

---

### `platosl info`

Show detailed information about a CUE schema.

```bash
platosl info <schema> [flags]

Flags:
      --format string    Output format (text, json, yaml) (default "text")
```

**Examples:**
```bash
# Text format (default)
platosl info schemas/person.cue

# JSON format
platosl info schemas/person.cue --format json

# YAML format
platosl info schemas/person.cue --format yaml
```

---

## Configuration File (platosl.yaml)

```yaml
version: v1
name: my-project

# Schema import paths
imports:
  - platosl.org/base/address/us@v1
  - platosl.org/base/address/uk@v1
  - ./custom/schemas

# Directories to validate
schemas:
  - schemas/
  - content/

# Validation options
validation:
  strict: true
  failOnWarning: false

# Code generation targets
generate:
  typescript:
    enabled: true
    output: generated/types.ts
    options:
      zod: true

  jsonschema:
    enabled: true
    output: generated/schema.json

  go:
    enabled: true
    output: generated/types.go
    options:
      package: types

  elixir:
    enabled: true
    output: generated/types.ex
    options:
      module: MyApp.Types
```

## Global Flags

Available on all commands:

```bash
  --config string    Config file (default "platosl.yaml")
  -v, --verbose      Verbose output
```

## Type Mappings

### CUE to TypeScript

| CUE Type | TypeScript Type |
|----------|----------------|
| `string` | `string` |
| `int` | `number` |
| `number`, `float` | `number` |
| `bool` | `boolean` |
| `[...T]` | `T[]` |
| `{field!: T}` | `{ field: T }` |
| `{field?: T}` | `{ field?: T }` |

### CUE to Go

| CUE Type | Go Type |
|----------|---------|
| `string` | `string` |
| `int` | `int` |
| `number`, `float` | `float64` |
| `bool` | `bool` |
| `[...T]` | `[]T` |
| `{field!: T}` | `Field T \`json:"field"\`` |
| `{field?: T}` | `Field *T \`json:"field,omitempty"\`` |

### CUE to Elixir

| CUE Type | Elixir Type |
|----------|-------------|
| `string` | `String.t()` |
| `int` | `integer()` |
| `number`, `float` | `float()` |
| `bool` | `boolean()` |
| `[...T]` | `list(T)` |
| `{field?: T}` | `T \| nil` |

## Examples

### Example 1: Blog Schema with Multiple Languages

**schemas/blog.cue:**
```cue
package schemas

#Post: {
	title!: string
	slug!:  string & =~"^[a-z0-9-]+$"
	body!:  string
	author!: #Author
	tags?: [...string]
	publishedAt?: string
}

#Author: {
	name!:  string
	email!: string & =~"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	bio?:   string
}
```

**Generate for all targets:**
```bash
platosl build
```

### Example 2: E-commerce Product Schema

**schemas/product.cue:**
```cue
package schemas

#Product: {
	id!:          string
	name!:        string
	description!: string
	price!:       number & >0
	currency!:    "USD" | "EUR" | "GBP"
	inStock!:     bool
	images?: [...string]
	categories?: [...string]
}
```

**Generate TypeScript with Zod:**
```bash
platosl gen typescript --zod --output src/types/product.ts
```

### Example 3: CI/CD Integration

**.github/workflows/validate.yml:**
```yaml
name: Validate Schemas

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install PlatoSL
        run: go install github.com/platoorg/platosl-cli@latest

      - name: Check formatting
        run: platosl fmt --check

      - name: Validate schemas
        run: platosl validate --strict

      - name: Generate code
        run: platosl build
```

## Troubleshooting

### "config file not found"

Run `platosl init` to create a new configuration file.

### "'cue' command not found" (for `platosl fmt`)

Install the CUE CLI:
```bash
# macOS
brew install cue-lang/tap/cue

# Linux/Other
go install cuelang.org/go/cmd/cue@latest
```

### "no schema paths configured"

Add schema paths to your `platosl.yaml`:
```yaml
schemas:
  - schemas/
```

### Validation errors

Use `--verbose` for more detailed output:
```bash
platosl validate --verbose
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

See [LICENSE](LICENSE) for details.

# PlatoSL CLI

A command-line tool for managing CUE-based schemas with validation and code generation for TypeScript, Zod, Go, JSON Schema, and Elixir.

**Key Features:**
- Build-time schema validation using CUE
- Import and extend battle-tested base schemas
- Generate type-safe code for multiple languages
- Framework-agnostic with support for any CMS or API

---

## Installation

### Prerequisites

- Go 1.24+ ([download here](https://go.dev/dl/))

### Install via Go

```bash
go install github.com/platoorg/platosl-cli/cmd/platosl@latest
```

### Install from Source

```bash
git clone https://github.com/platoorg/platosl-cli.git
cd platosl-cli
make install
```

### Verify Installation

```bash
platosl version
```

---

## Quick Start

### 1. Initialize a New Project

```bash
mkdir my-project
cd my-project
platosl init
```

You'll be prompted to select generators (use space to select, enter to confirm):
```
? Select generators to enable:
  [x] typescript
  [x] zod
  [ ] go
  [ ] jsonschema
  [ ] elixir
```

This creates:
- `platosl.yaml` - Configuration file
- `schemas/` - Your schema directory
- `schemas/example.cue` - Example schema
- `generated/` - Output directory for generated code

### 2. Define Your Schema

Edit `schemas/example.cue`:

```cue
package schemas

// Define a Person schema
#Person: {
	name!: string
	email!: string & =~"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	age?: int & >=0 & <=150
}
```

### 3. Validate Your Schema

```bash
platosl validate
```

### 4. Generate Code

```bash
# Generate TypeScript types
platosl gen typescript

# Generate Zod schemas
platosl gen zod

# Generate all configured generators
platosl build
```

Your generated code will be in the `generated/` directory.

---

## CLI Commands

| Command | Description |
|---------|-------------|
| `platosl init` | Initialize a new project with interactive generator selection |
| `platosl validate` | Validate all schemas in your project |
| `platosl gen <generator>` | Generate code for a specific generator |
| `platosl build` | Validate schemas and run all enabled generators |
| `platosl fmt` | Format CUE schema files |
| `platosl version` | Show version information |
| `platosl completion <shell>` | Generate shell completion script |

### Examples

```bash
# Initialize with specific generators (non-interactive)
platosl init --generators typescript,zod,go

# Validate a specific file
platosl validate schemas/person.cue

# Generate TypeScript types
platosl gen typescript

# Generate Zod schemas
platosl gen zod

# Validate and generate all enabled generators
platosl build

# Format all CUE files
platosl fmt
```

---

## Supported Frameworks & Languages

PlatoSL generates type-safe code for multiple languages and validation frameworks:

### TypeScript
Generates TypeScript interface definitions.

```typescript
export interface Person {
  name: string;
  email: string;
  age?: number;
}
```

**Configuration:**
```yaml
generate:
  typescript:
    enabled: true
    output: generated/types.ts
```

### Zod
Generates Zod validation schemas with TypeScript type inference.

```typescript
import { z } from 'zod';

export const PersonSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  age: z.number().min(0).max(150).optional(),
});

export type Person = z.infer<typeof PersonSchema>;
```

**Configuration:**
```yaml
generate:
  zod:
    enabled: true
    output: generated/schemas.ts
```

### Go
Generates Go struct definitions with JSON tags.

```go
package types

type Person struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   *int   `json:"age,omitempty"`
}
```

**Configuration:**
```yaml
generate:
  go:
    enabled: true
    output: generated/types.go
    options:
      package: types
```

### JSON Schema
Generates JSON Schema (Draft 2020-12) definitions.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "email": { "type": "string", "format": "email" },
    "age": { "type": "integer", "minimum": 0, "maximum": 150 }
  },
  "required": ["name", "email"]
}
```

**Configuration:**
```yaml
generate:
  jsonschema:
    enabled: true
    output: generated/schema.json
```

### Elixir
Generates Elixir typespecs and struct definitions.

```elixir
defmodule MyApp.Types do
  @type person :: %{
    name: String.t(),
    email: String.t(),
    age: integer() | nil
  }

  defstruct [:name, :email, :age]
end
```

**Configuration:**
```yaml
generate:
  elixir:
    enabled: true
    output: generated/types.ex
    options:
      module: MyApp.Types
```

---

## Adding Plato Schemas

PlatoSL provides battle-tested base schemas for common types like addresses, geographic data, and content blocks.

### Available Schema Categories

Official schemas are maintained at [github.com/platoorg/plato-sl](https://github.com/platoorg/plato-sl):

- **Addresses:** US, UK, DE, JP, and more country-specific address formats
- **Geographic data:** States, prefectures, regions, postal codes
- **Content blocks:** Image, Avatar, Card, Hero, and more CMS components

### Adding a Schema

Add schemas to your `platosl.yaml` imports section:

```bash
# Example: Add US address schema
platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0
```

Or edit `platosl.yaml` manually:

```yaml
imports:
  - github.com/platoorg/plato-sl/base/address/us@v1.0.0
  - github.com/platoorg/plato-sl/base/content@v1.0.0
```

### Using an Imported Schema

```cue
package myproject

import us "platosl.org/schemas/address/us"

// Use the base schema as-is
Address: us.#Address

// Or extend it with custom fields
MyAddress: us.#Address & {
    delivery_notes?: string

    // Add custom validation
    if state == "CA" {
        county!: string
    }
}
```

### Multi-Country Support

```cue
import (
    us "platosl.org/schemas/address/us"
    uk "platosl.org/schemas/address/uk"
    de "platosl.org/schemas/address/de"
)

// Union type for multi-country addresses
Address: us.#Address | uk.#Address | de.#Address
```

### Browse Available Schemas

Visit the [PlatoSL schemas catalog](https://github.com/platoorg/plato-sl/blob/main/base/README.md) to browse all available schemas.

---

## Configuration Reference

### platosl.yaml Structure

```yaml
version: v1
name: my-project

# Schema directories
schemas:
  - schemas/
  - custom-schemas/

# Schema dependencies
imports:
  - github.com/platoorg/plato-sl/base/address/us@v1.0.0
  - github.com/platoorg/plato-sl/base/content@v1.0.0

# Validation rules
validation:
  strict: true           # Require all fields to be concrete
  failOnWarning: false   # Don't fail on warnings

# Code generation
generate:
  typescript:
    enabled: true
    output: generated/types.ts
    options: {}

  zod:
    enabled: true
    output: generated/schemas.ts
    options: {}

  go:
    enabled: true
    output: generated/types.go
    options:
      package: types

  jsonschema:
    enabled: true
    output: generated/schema.json
    options: {}

  elixir:
    enabled: true
    output: generated/types.ex
    options:
      module: MyApp.Types
```

### Configuration Options

#### version (required)
Configuration file format version. Currently only `v1` is supported.

#### name (required)
Project name used in generated documentation and code comments.

#### schemas (required)
List of directories containing CUE schema files. Paths are relative to project root.

#### imports (optional)
List of base schema dependencies from the official PlatoSL repository.

#### validation (optional)
- `strict` - Require all fields to be concrete (fully defined)
- `failOnWarning` - Fail validation if warnings are encountered

#### generate (optional)
Configure code generators. Each generator has:
- `enabled` - Enable/disable the generator
- `output` - Output file path
- `options` - Generator-specific options (see Supported Frameworks section)

---

## Shell Completion

Enable autocompletion for faster command entry.

### Bash

```bash
# Load for current session
source <(platosl completion bash)

# Install permanently
# Linux:
sudo platosl completion bash > /etc/bash_completion.d/platosl

# macOS:
platosl completion bash > $(brew --prefix)/etc/bash_completion.d/platosl
```

### Zsh

```bash
# Enable completion
mkdir -p ~/.zsh/completions
platosl completion zsh > ~/.zsh/completions/_platosl

# Add to ~/.zshrc
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit; compinit' >> ~/.zshrc

# Restart shell
source ~/.zshrc
```

### Fish

```bash
# Load for current session
platosl completion fish | source

# Install permanently
platosl completion fish > ~/.config/fish/completions/platosl.fish
```

### PowerShell

```powershell
# Load for current session
platosl completion powershell | Out-String | Invoke-Expression

# Add to profile
platosl completion powershell > platosl.ps1
# Then source this file from your PowerShell profile
```

---

## Project Structure

A typical PlatoSL project:

```
my-project/
├── platosl.yaml          # Configuration file
├── schemas/              # Your CUE schemas
│   ├── person.cue
│   ├── product.cue
│   └── address.cue
└── generated/            # Generated code (auto-created)
    ├── types.ts          # TypeScript interfaces
    └── schemas.ts        # Zod schemas
```

---

## CUE Schema Basics

PlatoSL uses CUE (Configure Unify Execute) for schema definition. Here are the essentials:

### Define a Schema

```cue
package schemas

// Use # prefix for definitions
#Person: {
    name!: string              // Required field
    email?: string             // Optional field
    age: int & >=0 & <=150    // Integer with constraints
}
```

### Validation Rules

```cue
// String patterns
zip: string & =~"^\d{5}(-\d{4})?$"

// Number constraints
price: number & >0 & <1000000

// Enums
status: "active" | "inactive" | "pending"

// Nested objects
address: {
    street: string
    city: string
    zip: string
}

// Arrays
tags: [...string]  // Array of strings
```

### Composition

```cue
// Base schema
#BaseUser: {
    id: string
    name: string
}

// Extended schema
#AdminUser: #BaseUser & {
    permissions: [...string]
    role: "admin" | "superadmin"
}
```

Learn more about CUE at [cuelang.org](https://cuelang.org).

---

## Resources

- **Official Schemas:** [github.com/platoorg/plato-sl](https://github.com/platoorg/plato-sl)
- **Schema Catalog:** [Browse available schemas](https://github.com/platoorg/plato-sl/blob/main/base/README.md)
- **CUE Documentation:** [cuelang.org](https://cuelang.org)
- **Issue Tracker:** [github.com/platoorg/platosl-cli/issues](https://github.com/platoorg/platosl-cli/issues)

---

## License

[Add your license here]

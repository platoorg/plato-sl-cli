# PlatoSL - Schema Language for Content

**Kustomize for content. Strong validation. Best practices built-in.**

Version: 0.1 (Draft)

---

## What is PlatoSL?

PlatoSL is a **schema language with composition and validation** inspired by Kustomize and Kubernetes. It provides:

1. **Base schemas** - Battle-tested content structures (addresses, products, articles)
2. **Strong validation** - Catch errors at build time, not runtime
3. **Composition** - Extend and customize base schemas
4. **Framework agnostic** - Use with any CMS, UI framework, or API

---

## The Problem

Every project recreates the wheel:

```typescript
// ❌ Everyone writes their own address schema
interface Address {
  street?: string;  // Is this required?
  city: string;     // What about UK "town" vs US "city"?
  zip: string;      // Wrong: UK uses "postcode", format varies by country
}
```

Problems:
- No validation (runtime errors)
- Country-specific rules missed
- Inconsistent across projects
- No best practices

---

## The PlatoSL Solution

Use battle-tested base schemas with strong validation:

```bash
# Initialize project
platosl init

# Add base schemas from the schemas repository
platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0

# Validate at build time
platosl validate

# Generate TypeScript types
platosl gen typescript
```

---

## Quick Example: Addresses

### Use a Base Schema (US Address)

```cue
// my-project/address.cue
package myproject

import us "platosl.org/schemas/address/us"

// Use the US address base
Address: us.#Address

// Validation happens automatically
validAddress: Address & {
    street_line1: "123 Main St"
    city: "San Francisco"
    state: "CA"
    zip: "94102"
}

// ❌ This would fail validation:
// invalidAddress: Address & {
//     zip: "1234"  // Error: must match ^\d{5}(-\d{4})?$
// }
```

### Extend a Base Schema

```cue
// Add custom fields while keeping validation
MyAddress: us.#Address & {
    // Add optional delivery instructions
    delivery_notes?: string

    // Add custom validation
    if state == "CA" {
        // California requires county
        county!: string
    }
}
```

### Support Multiple Countries

```cue
package myproject

import (
    us "platosl.org/schemas/address/us"
    uk "platosl.org/schemas/address/uk"
    de "platosl.org/schemas/address/de"
)

// Union type for multi-country addresses
Address: us.#Address | uk.#Address | de.#Address

// Type-safe usage
myAddresses: [...Address] & [
    {
        _type: "us"
        street_line1: "123 Main St"
        city: "New York"
        state: "NY"
        zip: "10001"
    },
    {
        _type: "uk"
        street_line1: "10 Downing Street"
        town: "London"
        postcode: "SW1A 2AA"
    },
]
```

---

## Core Concepts

### 1. Base Schemas (Best Practices)

PlatoSL provides base schemas for common content types through the [platosl repository](https://github.com/platoorg/plato-sl):

```
platosl.org/schemas/
├── address/
│   ├── us/       # US addresses
│   ├── uk/       # UK addresses
│   ├── de/       # German addresses
│   ├── jp/       # Japanese addresses
│   └── ...
├── geo/          # Geographic data
│   ├── us/       # US states
│   ├── uk/       # UK countries
│   ├── de/       # German Bundesländer
│   └── jp/       # Japanese prefectures
└── content/      # CMS content blocks (Image, Avatar, etc.)
```

Add schemas to your project:
```bash
platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0
platosl add github.com/platoorg/plato-sl/base/content@v1.0.0
```

### 2. Strong Validation

Validation happens at **build time**, not runtime:

```bash
# Validate all schemas
platosl validate

# Output:
# ✓ address.cue validated
# ✗ product.cue failed:
#   - price.amount: must be positive
#   - image.src: required field missing
```

### 3. Composition & Extension

Like Kustomize, you can:
- **Use bases as-is**
- **Extend with new fields**
- **Override constraints**
- **Merge multiple bases**

```cue
// Start with base schema (after adding it to your project)
import content "platosl.org/schemas/content"

// Extend for your needs
MyImage: content.#Image & {
    // Add custom field
    internal_id?: string

    // Add constraints
    width: >800  // Minimum width requirement
}
```

### 4. Code Generation

Generate types for any language:

```bash
# TypeScript
platosl gen typescript > types.ts

# Go
platosl gen go > types.go

# JSON Schema
platosl gen jsonschema > schema.json

# GraphQL
platosl gen graphql > schema.graphql
```

---

## Installation

### Prerequisites

PlatoSL requires Go 1.24+ to be installed on your system.

### Install via Go

The easiest way to install PlatoSL is using `go install`:

```bash
go install github.com/platoorg/platosl-cli/cmd/platosl@latest
```

### Install from Source

Alternatively, you can build and install from source:

```bash
# Clone the repository
git clone https://github.com/platoorg/platosl-cli.git
cd platosl-cli

# Build and install
go install ./cmd/platosl

# Verify installation
platosl --help
```

### Install CUE (Optional)

While PlatoSL includes CUE internally, you may want to install the CUE CLI for direct schema manipulation:

```bash
go install cuelang.org/go/cmd/cue@latest
```

### Build with Version Information

If you're building from source and want version information embedded:

```bash
# Using make (recommended)
make build

# Or using go install with version flags
go install -ldflags "-X github.com/platoorg/platosl-cli/internal/cli.Version=v0.1.0" ./cmd/platosl
```

### Verify Installation

```bash
# Check PlatoSL version
platosl version

# Or use the --version flag
platosl --version

# Initialize a test project
mkdir test-project
cd test-project
platosl init
```

---

## Project Structure

```
my-project/
├── platosl.yaml          # Project config
├── schemas/
│   ├── address.cue       # Address schema
│   ├── product.cue       # Product schema
│   └── article.cue       # Article schema
├── content/
│   ├── products/         # Product content instances
│   └── articles/         # Article content instances
└── generated/
    ├── types.ts          # Generated TypeScript
    └── schema.json       # Generated JSON Schema
```

### platosl.yaml Configuration

The `platosl.yaml` file is the main configuration file for your PlatoSL project. It defines schema locations, dependencies, validation rules, and code generation settings.

#### Basic Example

```yaml
version: v1
name: my-project

# Schema locations
schemas:
  - schemas/

# Schema dependencies (added via `platosl add`)
imports:
  - github.com/platoorg/plato-sl/base/address/us@v1.0.0
  - github.com/platoorg/plato-sl/base/content@v1.0.0

# Validation rules
validation:
  strict: true
  failOnWarning: false

# Code generation
generate:
  typescript:
    enabled: true
    output: generated/types.ts
  jsonschema:
    enabled: true
    output: generated/schema.json
```

#### Configuration Reference

##### version (required)
The configuration file format version. Currently only `v1` is supported.

```yaml
version: v1
```

##### name (required)
The name of your project. Used in generated documentation and code comments.

```yaml
name: my-project
```

##### schemas (required)
List of directories containing your CUE schema files. Relative paths are resolved from the project root.

```yaml
schemas:
  - schemas/
  - custom-schemas/
```

##### imports (optional)
List of base schema dependencies to import into your project. These are typically versioned schema packages from the official PlatoSL repository.

```yaml
imports:
  - github.com/platoorg/plato-sl/base/address/us@v1.0.0
  - github.com/platoorg/plato-sl/base/address/uk@v1.0.0
  - github.com/platoorg/plato-sl/base/content@v1.0.0
```

##### validation (optional)
Configure schema validation behavior.

```yaml
validation:
  # Require all fields to be concrete (fully defined)
  strict: true

  # Fail the validation if warnings are encountered
  failOnWarning: false
```

##### generate (optional)
Configure code generation for different target languages. Each generator can be enabled/disabled and configured independently.

#### Available Generators

##### TypeScript Generator

Generates TypeScript interface definitions from your CUE schemas.

```yaml
generate:
  typescript:
    enabled: true
    output: generated/types.ts
    options: {}
```

**Output Example:**
```typescript
export interface Person {
  name: string;
  email: string;
  age?: number;
}
```

##### Zod Generator

Generates Zod validation schemas with TypeScript type inference.

```yaml
generate:
  zod:
    enabled: true
    output: generated/schemas.ts
    options: {}
```

**Output Example:**
```typescript
import { z } from 'zod';

export const PersonSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  age: z.number().min(0).max(150).optional(),
});

export type Person = z.infer<typeof PersonSchema>;
```

##### Go Generator

Generates Go struct definitions with JSON tags.

```yaml
generate:
  go:
    enabled: true
    output: generated/types.go
    options:
      package: types  # Go package name
```

**Output Example:**
```go
package types

type Person struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   *int   `json:"age,omitempty"`
}
```

##### JSON Schema Generator

Generates JSON Schema (Draft 2020-12) definitions.

```yaml
generate:
  jsonschema:
    enabled: true
    output: generated/schema.json
    options: {}
```

**Output Example:**
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/person.schema.json",
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "email": { "type": "string", "format": "email" },
    "age": { "type": "integer", "minimum": 0, "maximum": 150 }
  },
  "required": ["name", "email"]
}
```

##### Elixir Generator

Generates Elixir typespecs and struct definitions.

```yaml
generate:
  elixir:
    enabled: true
    output: generated/types.ex
    options:
      module: MyApp.Types  # Elixir module name
```

**Output Example:**
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

#### Complete Example

Here's a complete `platosl.yaml` with all generators enabled:

```yaml
version: v1
name: my-ecommerce-project

schemas:
  - schemas/
  - custom-schemas/

imports:
  - github.com/platoorg/plato-sl/base/address/us@v1.0.0
  - github.com/platoorg/plato-sl/base/address/uk@v1.0.0
  - github.com/platoorg/plato-sl/base/content@v1.0.0

validation:
  strict: true
  failOnWarning: false

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
      package: models

  jsonschema:
    enabled: true
    output: generated/schema.json
    options: {}

  elixir:
    enabled: true
    output: generated/types.ex
    options:
      module: MyEcommerce.Schemas
```

#### Interactive Initialization

When you run `platosl init`, you'll be prompted to select which generators to enable:

```bash
$ platosl init

? Select generators to enable:
  [x] typescript
  [x] zod
  [ ] go
  [ ] jsonschema
  [ ] elixir
```

Use the space bar to select/deselect generators, then press Enter to confirm. You can also specify generators non-interactively:

```bash
# Skip interactive mode and specify generators directly
platosl init --generators typescript,zod,go
```

---

## CLI Commands

```bash
# Show version information
platosl version

# Initialize new project
platosl init

# Add schema dependencies
platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0

# Validate schemas
platosl validate

# Validate specific file
platosl validate schemas/address.cue

# Generate code
platosl gen typescript
platosl gen jsonschema
platosl gen graphql

# Build (validate + generate)
platosl build

# Format CUE files
platosl fmt

# Show schema info
platosl info address
```

---

## Why CUE?

**CUE** (Configure Unify Execute) is perfect for schemas:

### Built-in Validation
```cue
// Validation is part of the schema
zip: string & =~"^\d{5}(-\d{4})?$"
price: number & >0 & <1000000
```

### Type Safety
```cue
// Types are constraints
name: string
age: int & >=0 & <=150
```

### Composition
```cue
// Schemas compose naturally
Base: { name: string }
Extended: Base & { age: int }
```

### Generate Anything
```bash
# CUE can export to JSON, YAML, Go, etc.
cue export schema.cue --out json
cue export schema.cue --out yaml
```

### Used by Kubernetes
- Kubernetes is moving from YAML to CUE
- Istio uses CUE
- Proven at scale

---

## Benefits

### For Content Creators
- ✅ Validation catches errors immediately
- ✅ Don't need to reinvent schemas
- ✅ Confidence that content will work

### For Developers
- ✅ Type-safe props (TypeScript, Go, etc.)
- ✅ No runtime validation code needed
- ✅ Contract between content and code

### For Organizations
- ✅ Consistent schemas across projects
- ✅ Best practices built-in
- ✅ Country-specific rules handled
- ✅ Migrate between CMSs easily

---

## Use Cases

### 1. Headless CMS Configuration

```cue
// Define content model for Contentful/Strapi
// First: platosl add github.com/platoorg/plato-sl/base/content@v1.0.0
import content "platosl.org/schemas/content"

CMSArticle: {
    // Use base content blocks
    hero?: content.#Image
    author_avatar?: content.#Avatar

    // CMS-specific fields
    _contentType: "article"
    _sys: {
        id: string
        createdAt: string
        updatedAt: string
    }
}
```

### 2. Multi-Country E-commerce

```cue
// First: platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0
// First: platosl add github.com/platoorg/plato-sl/base/address/uk@v1.0.0
import (
    us "platosl.org/schemas/address/us"
    uk "platosl.org/schemas/address/uk"
)

// Support both US and UK addresses
ShippingAddress: us.#Address | uk.#Address

Order: {
    customer_address: ShippingAddress
    billing_address: ShippingAddress
    // Validation ensures addresses are valid for their country
}
```

### 3. UI Component Props

```cue
// Define component props
ProductCard: {
    image: {
        src: string
        alt: string
    }
    title: string
    price: {
        amount: number & >0
        currency: "USD" | "EUR" | "GBP"
    }
}

// Generate TypeScript
// platosl gen typescript > ProductCard.types.ts
```

---

## Roadmap

### v0.1 (Current)
- [x] Core concept and architecture
- [ ] Base address schemas (US, UK, DE, JP, etc.)
- [ ] CLI tool (`platosl`)
- [ ] CUE schema library

### v0.2
- [ ] More base schemas (product, article, video, profile)
- [ ] TypeScript code generation
- [ ] JSON Schema export
- [ ] Validation errors with suggestions

### v0.3
- [ ] CMS integrations (Contentful, Strapi)
- [ ] Framework adapters (React, Vue)
- [ ] GraphQL schema generation
- [ ] Online schema browser

### v1.0
- [ ] Stable specification
- [ ] Full country coverage (addresses)
- [ ] Community schema registry
- [ ] IDE plugins

---

## Comparison

### vs JSON Schema
- JSON Schema: Validation only
- PlatoSL: Validation + composition + best practices

### vs TypeScript
- TypeScript: Runtime types in JS/TS only
- PlatoSL: Build-time validation, language-agnostic

### vs Kustomize
- Kustomize: For Kubernetes resources
- PlatoSL: For content and data

### vs Zod/Yup
- Zod/Yup: Runtime validation in JavaScript
- PlatoSL: Build-time validation, any language

---

## Philosophy

**Content should be validated like code.**

Just as we:
- Type-check code at compile time
- Validate Kubernetes manifests with `kubectl apply --dry-run`
- Lint configurations before deployment

We should:
- Validate content at build time
- Use proven schemas, not ad-hoc types
- Catch errors before production

**PlatoSL = Kustomize for content**

---

## Official Schemas

Official schemas are maintained in a separate repository:
[github.com/platoorg/plato-sl](https://github.com/platoorg/plato-sl)

Available schemas:
- **Address schemas** (US, UK, DE, JP)
- **Geographic data** (states, prefectures, regions)
- **Content contracts** (Image, Avatar, Card, Hero, etc.)

See the [schemas catalog](https://github.com/platoorg/plato-sl/blob/main/base/README.md).

## Next Steps

1. Read the [Quick Start Guide](./QUICKSTART.md)
2. Explore [official schemas](https://github.com/platoorg/plato-sl)
3. Check out the [CLI documentation](./CLI.md)
4. Try the [examples](https://github.com/platoorg/plato-sl/tree/main/examples)

---

## Contributing

PlatoSL is open source. Contributions needed:

- [ ] More country-specific address schemas
- [ ] Base schemas for common content types
- [ ] Code generators (TypeScript, Go, Python)
- [ ] CMS integrations
- [ ] Documentation and examples

---

**PlatoSL: Validate content like code. Compose like Kustomize.**

## PlatoSL Quick Start

**"Kustomize for content" - Strong validation for schemas, built-in best practices**

---

## What We Built

PlatoSL is a **schema language with composition and validation** for content. It's like Kustomize but for content instead of Kubernetes resources.

### The Core Idea

```
Base Schemas          Your Project          Generated Code
(best practices)  +   (customization)   →   (types, validation)

┌─────────────┐      ┌──────────────┐      ┌────────────────┐
│  US Address │      │ Extend with  │      │  TypeScript    │
│  UK Address │  +   │ custom fields│  →   │  Go structs    │
│  DE Address │      │ Override     │      │  JSON Schema   │
│  JP Address │      │ validation   │      │  GraphQL       │
└─────────────┘      └──────────────┘      └────────────────┘

      CUE                  CUE                  Any language
   (validate)           (compose)              (use everywhere)
```

---

## What You Asked For

✅ **Schemas for addresses by country** - US, UK, DE, JP (more coming)
✅ **Strong validation** - Build-time checking like Kustomize
✅ **Best practices built-in** - USPS, Royal Mail, Deutsche Post standards
✅ **Headless CMS ready** - Use in Contentful, Strapi, etc.
✅ **Composable** - Extend/override like Kustomize overlays
✅ **Using CUE** - Better than YAML for validation

---

## Repository Structure

PlatoSL is now split into two repositories:

### platosl-cli (This Repository)
The command-line tool for working with schemas:
```
platosl-cli/
├── README.md                    # CLI documentation
├── QUICKSTART.md                # This file
├── cmd/platosl/                 # CLI tool source
├── internal/                    # CLI implementation
└── testdata/                    # Test schemas
```

### platosl (Schemas Repository)
Official base schemas maintained separately:
```
platosl/
├── base/                        # Base schemas (like Kustomize bases)
│   ├── address/
│   │   ├── us/address.cue      # US addresses (USPS standard)
│   │   ├── uk/address.cue      # UK addresses (Royal Mail)
│   │   ├── de/address.cue      # German addresses (Deutsche Post)
│   │   └── jp/address.cue      # Japanese addresses (Japan Post)
│   ├── geo/                     # Geographic data
│   │   ├── us/states.cue       # US states
│   │   ├── uk/countries.cue    # UK countries
│   │   ├── de/states.cue       # German Bundesländer
│   │   └── jp/prefectures.cue  # Japanese prefectures
│   └── content/                 # CMS content blocks
│       ├── image.cue           # Image schemas
│       ├── avatar.cue          # Avatar schemas
│       └── blocks.cue          # Content blocks
└── examples/                    # Usage examples
    ├── 01_basic_usage.cue      # Import and use base schemas
    ├── 02_extending_schemas.cue # Extend/customize
    └── 03_cms_integration.cue   # Use with CMS
```

---

## How It Works

### 1. Initialize Project and Add Schemas

```bash
# Initialize a new PlatoSL project
platosl init

# Add base schemas from the schemas repository
platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0
platosl add github.com/platoorg/plato-sl/base/content@v1.0.0
```

### 2. Use Base Schemas (Best Practices)

```cue
// Import a base schema (after adding it to your project)
import us "platosl.org/schemas/address/us"

// Use it
myAddress: us.#Address & {
    street_line1: "123 Main St"
    city: "San Francisco"
    state: "CA"
    zip: "94102"
}

// ✅ Validated automatically!
// ❌ Invalid data fails at build time
```

### 3. Extend for Your Needs (Like Kustomize Overlays)

```cue
// Add custom fields
MyAddress: us.#Address & {
    delivery_notes?: string
    gate_code?: string
}

// Override validation
CaliforniaOnly: us.#Address & {
    state: "CA"  // Only allow California
}
```

### 4. Generate Code for Any Language

```bash
platosl gen typescript  # → types.ts
platosl gen go          # → types.go
platosl gen jsonschema  # → schema.json
platosl gen graphql     # → schema.graphql
```

### 5. Use Everywhere

**Frontend (React/TypeScript):**
```tsx
import { USAddress } from './generated/types';
function AddressForm(props: { address: USAddress }) { ... }
```

**Backend (Go):**
```go
import "myproject/generated"
func HandleAddress(addr address.USAddress) { ... }
```

**CMS (Contentful):**
```javascript
// Import generated Contentful schema
// Content validated before it enters CMS
```

---

## Key Examples

### Example 1: Multi-Country E-commerce

```cue
// Support US, UK, Germany
// First add: platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0
// First add: platosl add github.com/platoorg/plato-sl/base/address/uk@v1.0.0
// First add: platosl add github.com/platoorg/plato-sl/base/address/de@v1.0.0
import (
    us "platosl.org/schemas/address/us"
    uk "platosl.org/schemas/address/uk"
    de "platosl.org/schemas/address/de"
)

Address: us.#Address | uk.#Address | de.#Address

customers: [...{
    name: string
    shipping: Address
    billing: Address
}]

// ✅ All addresses validated against country-specific rules
```

### Example 2: CMS Content Type

```cue
#StoreLocation: {
    store_name: string
    address: us.#Address  // Use validated schema
    phone: string
    hours: string
}

// Export to Contentful/Strapi
// platosl gen contentful > contentful-schema.json
```

### Example 3: California-Only Validation

```cue
CaliforniaAddress: us.#Address & {
    state: "CA"          // Must be California
    county!: string      // Require county
}

// ✅ Forces business rules at schema level
```

---

## Validation Examples

### ✅ Valid US Address

```cue
valid: us.#Address & {
    street_line1: "1600 Pennsylvania Ave NW"
    city: "Washington"
    state: "DC"
    zip: "20500"
}
// ✅ Passes validation
```

### ❌ Invalid US Address

```cue
invalid: us.#Address & {
    street_line1: "123 Main St"
    city: "New York"
    state: "NY"
    zip: "INVALID"  // ❌ Error: does not match ^\d{5}(-\d{4})?$
}
// Build fails with clear error message
```

---

## Benefits

### vs Writing Your Own Schemas

| Without PlatoSL | With PlatoSL |
|-----------------|--------------|
| ❌ Reinvent the wheel | ✅ Use battle-tested schemas |
| ❌ Miss country-specific rules | ✅ Built-in best practices |
| ❌ Runtime errors | ✅ Build-time validation |
| ❌ Manual sync across codebases | ✅ Generate for all languages |
| ❌ Drift between frontend/backend | ✅ Single source of truth |

### The PlatoSL Way

```
Define Once → Validate Early → Use Everywhere
   (CUE)      (Build time)     (All languages)
```

---

## Next Steps

1. **Install CUE** (the validation engine)
   ```bash
   go install cuelang.org/go/cmd/cue@latest
   ```

2. **Install PlatoSL CLI**
   ```bash
   go install github.com/platoorg/plato-sl-cli/cmd/platosl@latest
   ```

3. **Initialize a project**
   ```bash
   mkdir my-project && cd my-project
   platosl init
   ```

4. **Add base schemas**
   ```bash
   platosl add github.com/platoorg/plato-sl/base/address/us@v1.0.0
   ```

5. **Create your first schema**
   ```bash
   cat > schemas/user.cue <<EOF
   package schemas

   import us "platosl.org/schemas/address/us"

   #User: {
       name!: string
       address!: us.#Address
   }
   EOF
   ```

6. **Validate**
   ```bash
   platosl validate
   ```

7. **Explore official schemas**
   Visit [github.com/platoorg/plato-sl](https://github.com/platoorg/plato-sl) to see all available base schemas and examples.

---

## The Vision

### Today (v0.1)
- Base address schemas (US, UK, DE, JP)
- CUE validation
- Examples and docs

### Soon (v0.2)
- More countries
- CLI tool (`platosl validate`, `platosl gen`)
- TypeScript/Go code generation
- CMS adapters

### Future (v1.0)
- 50+ country address schemas
- Product, article, video schemas
- Framework integrations (React, Vue)
- Community schema registry
- IDE plugins

---

## Questions?

**Q: Do I need to use all countries?**
A: No, import only what you need.

**Q: Can I use this with my existing CMS?**
A: Yes! It's a layer on top of your CMS.

**Q: Is CUE hard to learn?**
A: It's like JSON with types. Check out the examples!

**Q: Can I extend base schemas?**
A: Absolutely! That's the whole point (like Kustomize).

**Q: Is this production-ready?**
A: v0.1 is draft/experimental. Use for prototyping.

---

**PlatoSL: Validate content like code. Compose like Kustomize.**

Built with ❤️ using [CUE](https://cuelang.org/)

#!/bin/bash
set -e

# Test script for PlatoSL CLI
echo "================================"
echo "Testing PlatoSL CLI"
echo "================================"
echo ""

# Cleanup
TEST_DIR="/tmp/platosl-test-$$"
rm -rf "$TEST_DIR"

echo "Test directory: $TEST_DIR"
echo ""

# Build the CLI
echo "Building CLI..."
make build
BIN="$(pwd)/bin/platosl"
echo "Binary: $BIN"
echo ""

# Test 1: Initialize project
echo "Test 1: platosl init"
echo "--------------------"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

$BIN init test-project --name "Test Project"
cd test-project

if [ -f "platosl.yaml" ]; then
    echo "✓ platosl.yaml created"
else
    echo "✗ platosl.yaml not found"
    exit 1
fi

if [ -d "schemas" ]; then
    echo "✓ schemas/ directory created"
else
    echo "✗ schemas/ directory not found"
    exit 1
fi

if [ -f "schemas/example.cue" ]; then
    echo "✓ schemas/example.cue created"
else
    echo "✗ schemas/example.cue not found"
    exit 1
fi
echo ""

# Test 2: Validate schemas
echo "Test 2: platosl validate"
echo "------------------------"
$BIN validate
echo "✓ Validation passed"
echo ""

# Test 3: Generate TypeScript
echo "Test 3: platosl gen typescript"
echo "-------------------------------"
$BIN gen typescript --output generated/types.ts
if [ -f "generated/types.ts" ]; then
    echo "✓ TypeScript generated"
    echo ""
    echo "Generated TypeScript (first 20 lines):"
    head -20 generated/types.ts
else
    echo "✗ TypeScript not generated"
    exit 1
fi
echo ""

# Test 4: Generate TypeScript with Zod
echo "Test 4: platosl gen typescript --zod"
echo "-------------------------------------"
$BIN gen typescript --zod --output generated/types-zod.ts
if [ -f "generated/types-zod.ts" ]; then
    echo "✓ TypeScript with Zod generated"
    if grep -q "import { z } from 'zod'" generated/types-zod.ts; then
        echo "✓ Zod import found"
    else
        echo "✗ Zod import not found"
        exit 1
    fi
else
    echo "✗ TypeScript with Zod not generated"
    exit 1
fi
echo ""

# Test 5: Generate JSON Schema
echo "Test 5: platosl gen jsonschema"
echo "-------------------------------"
$BIN gen jsonschema --output generated/schema.json
if [ -f "generated/schema.json" ]; then
    echo "✓ JSON Schema generated"
    echo ""
    echo "Generated JSON Schema (first 15 lines):"
    head -15 generated/schema.json
else
    echo "✗ JSON Schema not generated"
    exit 1
fi
echo ""

# Test 6: Generate Go
echo "Test 6: platosl gen go"
echo "----------------------"
$BIN gen go --output generated/types.go --package types
if [ -f "generated/types.go" ]; then
    echo "✓ Go code generated"
    echo ""
    echo "Generated Go code (first 20 lines):"
    head -20 generated/types.go
else
    echo "✗ Go code not generated"
    exit 1
fi
echo ""

# Test 7: Generate Elixir
echo "Test 7: platosl gen elixir"
echo "--------------------------"
$BIN gen elixir --output generated/types.ex --module TestProject.Types
if [ -f "generated/types.ex" ]; then
    echo "✓ Elixir code generated"
    echo ""
    echo "Generated Elixir code (first 20 lines):"
    head -20 generated/types.ex
else
    echo "✗ Elixir code not generated"
    exit 1
fi
echo ""

# Test 8: Info command
echo "Test 8: platosl info"
echo "--------------------"
$BIN info schemas/example.cue
echo "✓ Info command works"
echo ""

# Test 9: Build command
echo "Test 9: platosl build"
echo "---------------------"
$BIN build
echo "✓ Build command works"
echo ""

# Cleanup
echo "Cleaning up..."
cd /
rm -rf "$TEST_DIR"

echo ""
echo "================================"
echo "All tests passed! ✓"
echo "================================"

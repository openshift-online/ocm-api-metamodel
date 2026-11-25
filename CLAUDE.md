# CLAUDE.md

<!-- Canonical source: AGENTS.md. This file is auto-generated for Claude Code compatibility. -->

This file provides guidance to AI coding assistants when working with this repository.

## Project Overview

OCM API Metamodel Tool (`ocm-metamodel-tool`) — a code generator that processes OCM API model definitions and generates Go source code, documentation, and OpenAPI specifications for the OpenShift Cluster Manager API.

## Build & Test Commands

```bash
make binary          # Build the ocm-metamodel-tool binary
make generate        # Regenerate ANTLR parser code
make go_tests        # Run Go unit tests
make docs_tests      # Run documentation tests
make fmt             # Format Go source code
make clean           # Remove build artifacts
```

## Architecture

- **cmd/**: CLI entry point for `ocm-metamodel-tool`
- **pkg/**: Core packages
  - **pkg/language/**: ANTLR-based model definition parser
  - **pkg/concepts/**: In-memory model representation
  - **pkg/generators/**: Code generators for Go, docs, and OpenAPI
  - **pkg/nomenclator/**: Naming convention utilities
  - **pkg/reporter/**: Error and warning reporting
- **tests/**: Ginkgo-based test suites

## Key Conventions

- Uses ANTLR v4 for parsing model definitions
- Ginkgo/Gomega for testing
- Module path: `github.com/openshift-online/ocm-api-metamodel`
- Generated code follows Go conventions with proper type safety

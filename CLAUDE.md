# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the OpenShift Cluster Manager (OCM) API metamodel project containing the `ocm-metamodel-tool` command line utility. The tool generates code from API models defined in a custom domain-specific language (DSL) using ANTLR grammars.

## Key Commands

### Building
```bash
make binary          # Build the metamodel binary
```

### Code Generation
```bash
# Generate ANTLR parser/lexer code from grammar files
make generate

# The metamodel tool itself generates various outputs:
./metamodel generate go --model=tests/model --base=<base-package> --output=<output-dir>
./metamodel generate openapi --model=tests/model --output=<output-dir>
./metamodel generate docs --model=tests/model --output=<output-dir>
```

### Testing
```bash
make test            # Run all tests (unit, go, openapi, docs)
make unit_tests      # Run unit tests only
make go_tests        # Run Go generation tests
make openapi_tests   # Run OpenAPI validation tests
make docs_tests      # Run documentation generation tests
```

### Code Quality
```bash
make verify          # Run go vet and gofmt verification
make fmt             # Format code with gofmt
```

### Cleanup
```bash
make clean           # Remove generated files and build artifacts
```

## Architecture Overview

### Core Components

- **Language Parser** (`pkg/language/`): ANTLR-based parser for the metamodel DSL
  - `ModelLexer.g4` and `ModelParser.g4` define the grammar
  - Generated Go parser code processes `.model` files

- **Concepts** (`pkg/concepts/`): Core data structures representing the metamodel
  - `Model`: Root container for services and versions
  - `Service`: Groups of related API versions
  - `Version`: Contains types, resources, and methods
  - `Type`: Enums, classes, structs, and errors
  - `Resource`: REST API resources with CRUD operations
  - `Method`: Individual API operations

- **Generators** (`pkg/generators/`): Code generation backends
  - `golang/`: Generates Go client libraries, types, builders
  - `openapi/`: Generates OpenAPI specifications
  - `docs/`: Generates documentation

- **Supporting Packages**:
  - `pkg/names/`: Name conversion utilities (camelCase, snake_case, etc.)
  - `pkg/annotations/`: Annotation processing for model elements
  - `pkg/reporter/`: Error reporting and diagnostics

### Model Files

Model files (`.model` extension) use a custom DSL to define API structures:
- Located in `tests/model/` for test cases
- Define classes, enums, resources, and their relationships
- Support annotations for additional metadata
- Reference external types using `@ref` annotations

### Code Generation Flow

1. Parse `.model` files using ANTLR grammar
2. Build internal concept representation
3. Apply transformations and validations
4. Generate output using target-specific generators
5. Validate generated code (for OpenAPI, run through validator)

## Development Notes

- Uses ANTLR 4.9.3 for parsing (downloaded automatically)
- Tests include validation of generated Go code compilation
- OpenAPI specs are validated using OpenAPI Generator CLI
- Ginkgo v2 is used for testing framework
- CGO is disabled to ensure static binary builds

## Test Structure

- `tests/model/`: Sample model definitions for testing
- `tests/go/`: Go code generation tests
- `tests/openapi/`: OpenAPI generation and validation tests
- `tests/docs/`: Documentation generation tests
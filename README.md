# OpenShift cluster manager API metamodel

This project contains the code for the `ocm-metamodel-tool` command line utility,
which is used to generate code from the API model.

## Installation

### Prerequisites

- Go 1.17+
- Java runtime (for ANTLR parser generation)

### Building

```bash
make binary           # Build the metamodel binary
```

The binary is placed in the project root directory.

## Usage

This tool is used as part of the OCM SDK and API model code generation pipeline.
It processes model definitions from [ocm-api-model](https://github.com/openshift-online/ocm-api-model)
and generates Go source code for [ocm-sdk-go](https://github.com/openshift-online/ocm-sdk-go).

See the ocm-api-model and ocm-sdk-go repositories for the full generation workflow.

## Development

### Building from Source

```bash
make binary           # Build the tool
make generate         # Regenerate ANTLR parser code
```

### Testing

```bash
make test             # Run all tests
make unit_tests       # Run unit tests only
make go_tests         # Run Go integration tests
make openapi_tests    # Run OpenAPI validation tests
make docs_tests       # Run documentation tests
```

### Contributing

1. Fork this repository
2. Make your changes
3. Run `make test` to verify
4. Submit a pull request

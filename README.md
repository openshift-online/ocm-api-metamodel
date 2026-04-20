# OpenShift cluster manager API metamodel

This project contains the code for the `ocm-metamodel-tool` command line utility,
which is used to generate code from the API model.

## Installation

### Prerequisites

- Go 1.17+
- Java runtime (for ANTLR parser generation)

### Building

```bash
make binary           # Build the ocm-metamodel-tool binary
```

The binary is placed in the project root directory.

## Usage

Generate Go code from an API model:

```bash
./ocm-metamodel-tool generate   --model=/path/to/ocm-api-model/model   --output=/path/to/output
```

Run `./ocm-metamodel-tool --help` for all available commands and options.

## Development

### Building from Source

```bash
make binary           # Build the tool
make generate         # Regenerate ANTLR parser code
```

### Testing

```bash
make go_tests         # Run Go unit tests
make docs_tests       # Run documentation tests
```

### Contributing

1. Fork this repository
2. Make your changes
3. Run `make go_tests` to verify
4. Submit a pull request

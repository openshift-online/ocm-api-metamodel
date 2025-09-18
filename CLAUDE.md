I've analyzed the OCM API metamodel codebase and created a comprehensive CLAUDE.md file. The analysis revealed:

## Key Findings:

**Project Type**: Go-based code generation tool for API client libraries and specifications

**Architecture**: 
- Model parser using ANTLR4 for custom DSL
- Core concepts package representing API elements  
- Pluggable generators for Go, OpenAPI, and docs output
- Comprehensive test suite with integration testing

**Development Commands** (from Makefile):
- `make binary` - Build the metamodel CLI
- `make test` - Run all test suites  
- `make verify` - Code quality checks
- `make generate` - Generate ANTLR parser

**Testing Strategy**:
- Unit tests with Ginkgo/Gomega
- Integration tests generating code from test models
- Separate validation for each generator type

The CLAUDE.md file provides future Claude Code instances with essential information about the codebase structure, common commands, testing approach, and code generation workflow without duplicating obvious details or including generic development practices.

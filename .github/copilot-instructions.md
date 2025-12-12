# GitHub Copilot Instructions for this repo

## Communication Style
- Be blunt
- Write short sentences
- Use bullet points
- Avoid useless filler phrases
- Avoid useless explanations
- Write like giving instructions to a machine not writing to a human

## Project Overview
- golang cli tool to create workflow documentation from comments in workflow files with @tagname value syntx
- Cross-platform: Windows and Linux
- Cross-architecture: AMD64 and ARM64

## Technology Stack
- Go 1.25.3
- github.com/spf13/cobra for CLI
- gopkg.in/yaml.v3 for YAML parsing
- Latest golangci-lint for linting

## Code Standards
- Follow Go best practices
- Use TDD approach
- Write tests before implementation
- Maintain test coverage above 70%
- Run linter before commits
- Keep functions small and focused
- Prefer explicit over implicit

## Project Structure
```
.
├── cmd/workflowdocgen/         # Main application entry
├── pkg/workflowdocgenr/       # Core substitution logic
├── scripts/               # Build and test scripts
├── .github/workflows/     # CI/CD pipelines
└── README.md
```

## CLI Flags
- `--yaml`: YAML file with values (required)
- `--file`: Input file with placeholders (optional, defaults to stdin)
- `--help`: Show help
- `--version`: Show version info

## Testing
- Unit tests in `*_test.go` files
- Integration tests in `scripts/test.sh` and `scripts/test.ps1`
- Run tests: `./scripts/test.sh` or `./scripts/test.ps1`
- Coverage report included

## Building
- Build all platforms: `./scripts/build.sh` or `./scripts/build.ps1`
- Output in `dist/` directory
- Version/commit/date embedded at build time

## CI/CD
- GitHub Actions workflows in `.github/workflows/`
- CI runs on push/PR
- Release on tag push (v*)
- Automated binary uploads

## Development Workflow
1. Write test
2. Implement feature
3. Run tests
4. Run linter
5. Build locally
6. Commit

## Copilot Agent Requirements
- **ALWAYS** run linter before finishing/completing a PR
- Linter command: `golangci-lint run --timeout=5m`
- Fix all linter errors before marking work complete
- Do not commit code with linter errors

## Version Management
- Use semantic versioning
- Tag releases with `v` prefix (v1.0.0)
- Embed version in binary via ldflags

## Dependencies
- Keep dependencies minimal
- Update regularly
- Pin major versions
- Check security advisories

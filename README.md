# github-workflow-doc

A Go CLI tool that helps to document your GitHub workflows based on "javadoc" like comments.

## Overview

This tool generates a `WORKFLOWS.md` file that documents all workflows in your `.github/workflows` directory by parsing special documentation comments.

## Supported Documentation Comments

Add these comments to your workflow YAML files:

- `# @workflow.name:` - Name of the workflow
- `# @workflow.description:` - Description of what the workflow does
- `# @workflow.owners:` - Team or person responsible (e.g., team-release)
- `# @workflow.tags:` - Tags for categorization (comma-separated)
- `# @workflow.params:` - Input parameters the workflow accepts
- `# @workflow.results:` - Output or results produced by the workflow
- `# @workflow.permissions:` - Required permissions for the workflow
- `# @workflow.requirements:` - Setup steps needed before using the workflow
- `# @job.description:` - Description of a specific job
- `# @step.description:` - Description of a specific step

## Installation

### Requirements

- Go 1.23 or higher

### Build from Source

```bash
go build -o bin/workflowdocgen ./cmd/workflowdocgen
```

## Usage

Run the tool from your repository root:

```bash
./bin/workflowdocgen
```

### Options

- `--workflows-dir` - Path to workflows directory (default: `.github/workflows`)
- `--output` - Output file path (default: `WORKFLOWS.md`)

### Example

```bash
./bin/workflowdocgen --workflows-dir .github/workflows --output WORKFLOWS.md
```

## Example Workflow Documentation

```yaml
# @workflow.name: CI Build and Test
# @workflow.description: Runs continuous integration tests on all pull requests
# @workflow.owners: team-release
# @workflow.tags: ci, testing, automation
# @workflow.params: branch name, commit SHA
# @workflow.results: test results, coverage report
# @workflow.permissions: read repository contents, write test results
# @workflow.requirements: GitHub Actions enabled, test dependencies installed

name: CI Build and Test

on:
  pull_request:
    branches: [ main ]
```

## Output

The tool generates a `WORKFLOWS.md` file containing:

1. A markdown table with columns: Workflow | Description | Owners | Tags | File
2. Detailed workflow information section with params, results, permissions, and requirements

## Development

### Project Structure

```
.
├── cmd/
│   └── workflowdocgen/     # CLI entry point
│       └── main.go
├── pkg/
│   └── workflowdocgen/     # Library logic
│       ├── parser.go       # YAML parsing and comment extraction
│       └── generator.go    # Markdown generation
└── .github/
    └── workflows/          # Example workflow files
```

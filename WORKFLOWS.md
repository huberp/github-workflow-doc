# Workflow Documentation

This document provides an overview of all GitHub workflows in this repository.

| Workflow | Description | Owners | Tags | File |
|----------|-------------|--------|------|------|
| CI | Run linting, tests, and multi-platform builds on push/PR | - | - | ci.yml |
| Go Dependency Submission | Submit Go module dependencies to GitHub's dependency graph for security analysis | - | - | go-dependency-submission.yml |
| Release | Build and publish releases with multi-platform binaries using GoReleaser | - | - | release.yml |

## Detailed Workflow Information

### CI

**Permissions:** contents:read

### Go Dependency Submission

**Permissions:** contents:write

### Release

**Permissions:** contents:write


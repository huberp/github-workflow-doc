# Workflow Documentation

This document provides an overview of all GitHub workflows in this repository.

| Workflow | Description | Owners | Tags | File |
|----------|-------------|--------|------|------|
| CI Build and Test | Runs continuous integration tests on all pull requests | team-release | ci, testing, automation | ci.yml |
| Release | Creates a new release and publishes artifacts | team-release | release, deployment | release.yml |

## Detailed Workflow Information

### CI Build and Test

**Parameters:** branch name, commit SHA

**Results:** test results, coverage report

**Permissions:** read repository contents, write test results

**Requirements:** GitHub Actions enabled, test dependencies installed

### Release

**Permissions:** write releases, read repository contents

**Requirements:** Release tag must follow semver format


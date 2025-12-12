package workflowdocgen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseWorkflowFile(t *testing.T) {
	// Create temp directory for test files
	tempDir := t.TempDir()

	t.Run("valid workflow with all fields", func(t *testing.T) {
		content := `# @workflow.name: CI Pipeline
# @workflow.description: Continuous integration workflow
# @workflow.owners: team-platform
# @workflow.tags: ci, testing
# @workflow.params: branch, environment
# @workflow.results: build artifacts
# @workflow.permissions: read:contents
# @workflow.requirements: Go 1.23+

name: CI
on: push
`
		filePath := filepath.Join(tempDir, "test1.yml")
		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := ParseWorkflowFile(filePath)
		if err != nil {
			t.Fatalf("ParseWorkflowFile failed: %v", err)
		}

		if doc.Name != "CI Pipeline" {
			t.Errorf("Expected name 'CI Pipeline', got '%s'", doc.Name)
		}
		if doc.Description != "Continuous integration workflow" {
			t.Errorf("Expected description 'Continuous integration workflow', got '%s'", doc.Description)
		}
		if doc.Owners != "team-platform" {
			t.Errorf("Expected owners 'team-platform', got '%s'", doc.Owners)
		}
		if doc.Tags != "ci, testing" {
			t.Errorf("Expected tags 'ci, testing', got '%s'", doc.Tags)
		}
		if doc.Params != "branch, environment" {
			t.Errorf("Expected params 'branch, environment', got '%s'", doc.Params)
		}
		if doc.Results != "build artifacts" {
			t.Errorf("Expected results 'build artifacts', got '%s'", doc.Results)
		}
		if doc.Permissions != "read:contents" {
			t.Errorf("Expected permissions 'read:contents', got '%s'", doc.Permissions)
		}
		if doc.Requirements != "Go 1.23+" {
			t.Errorf("Expected requirements 'Go 1.23+', got '%s'", doc.Requirements)
		}
		if doc.FileName != "test1.yml" {
			t.Errorf("Expected filename 'test1.yml', got '%s'", doc.FileName)
		}
		if doc.FilePath != filePath {
			t.Errorf("Expected filepath '%s', got '%s'", filePath, doc.FilePath)
		}
	})

	t.Run("workflow with partial fields", func(t *testing.T) {
		content := `# @workflow.name: Test Workflow
# @workflow.description: Simple test

name: Test
on: pull_request
`
		filePath := filepath.Join(tempDir, "test2.yml")
		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := ParseWorkflowFile(filePath)
		if err != nil {
			t.Fatalf("ParseWorkflowFile failed: %v", err)
		}

		if doc.Name != "Test Workflow" {
			t.Errorf("Expected name 'Test Workflow', got '%s'", doc.Name)
		}
		if doc.Description != "Simple test" {
			t.Errorf("Expected description 'Simple test', got '%s'", doc.Description)
		}
		if doc.Owners != "" {
			t.Errorf("Expected empty owners, got '%s'", doc.Owners)
		}
		if doc.Tags != "" {
			t.Errorf("Expected empty tags, got '%s'", doc.Tags)
		}
	})

	t.Run("workflow with no documentation comments", func(t *testing.T) {
		content := `name: Plain Workflow
on: push
jobs:
  build:
    runs-on: ubuntu-latest
`
		filePath := filepath.Join(tempDir, "test3.yml")
		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := ParseWorkflowFile(filePath)
		if err != nil {
			t.Fatalf("ParseWorkflowFile failed: %v", err)
		}

		if doc.Name != "" {
			t.Errorf("Expected empty name, got '%s'", doc.Name)
		}
		if doc.Description != "" {
			t.Errorf("Expected empty description, got '%s'", doc.Description)
		}
	})

	t.Run("workflow with malformed comments", func(t *testing.T) {
		content := `# @workflow.name CI Pipeline
# @ workflow.description: Bad format
# @workflow.tags: valid-tag
# @workflow.invalid: Should be ignored

name: Test
`
		filePath := filepath.Join(tempDir, "test4.yml")
		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := ParseWorkflowFile(filePath)
		if err != nil {
			t.Fatalf("ParseWorkflowFile failed: %v", err)
		}

		// Valid tag should be parsed
		if doc.Tags != "valid-tag" {
			t.Errorf("Expected tags 'valid-tag', got '%s'", doc.Tags)
		}
		// Name without colon after field name should not be parsed
		if doc.Name != "" {
			t.Errorf("Expected empty name (malformed - no colon), got '%s'", doc.Name)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test5.yml")
		if err := os.WriteFile(filePath, []byte(""), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := ParseWorkflowFile(filePath)
		if err != nil {
			t.Fatalf("ParseWorkflowFile failed: %v", err)
		}

		if doc.Name != "" {
			t.Errorf("Expected empty name for empty file, got '%s'", doc.Name)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := ParseWorkflowFile("/nonexistent/file.yml")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})

	t.Run("workflow with whitespace variations", func(t *testing.T) {
		content := `# @workflow.name: Spaces Workflow
# @workflow.description: Valid description
# @workflow.tags: valid

name: Test
`
		filePath := filepath.Join(tempDir, "test6.yml")
		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := ParseWorkflowFile(filePath)
		if err != nil {
			t.Fatalf("ParseWorkflowFile failed: %v", err)
		}

		if doc.Name != "Spaces Workflow" {
			t.Errorf("Expected name 'Spaces Workflow', got '%s'", doc.Name)
		}
		if doc.Description != "Valid description" {
			t.Errorf("Expected description 'Valid description', got '%s'", doc.Description)
		}
		if doc.Tags != "valid" {
			t.Errorf("Expected tags 'valid', got '%s'", doc.Tags)
		}
	})

	t.Run("workflow with extra spaces after hash", func(t *testing.T) {
		content := `#  @workflow.name: Should Not Parse
name: Test
`
		filePath := filepath.Join(tempDir, "test7.yml")
		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		doc, err := ParseWorkflowFile(filePath)
		if err != nil {
			t.Fatalf("ParseWorkflowFile failed: %v", err)
		}

		// Pre-check requires "# @" prefix (exactly 1 space), so this should not parse
		if doc.Name != "" {
			t.Errorf("Expected empty name (pre-check should fail with 2 spaces), got '%s'", doc.Name)
		}
	})
}

func TestParseWorkflowsDirectory(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("directory with multiple workflow files", func(t *testing.T) {
		// Create test files
		files := map[string]string{
			"ci.yml": `# @workflow.name: CI
# @workflow.description: CI workflow
name: CI`,
			"test.yaml": `# @workflow.name: Test
# @workflow.description: Test workflow
name: Test`,
			"deploy.yml": `# @workflow.name: Deploy
name: Deploy`,
		}

		for name, content := range files {
			filePath := filepath.Join(tempDir, name)
			if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
				t.Fatalf("Failed to create test file %s: %v", name, err)
			}
		}

		docs, err := ParseWorkflowsDirectory(tempDir)
		if err != nil {
			t.Fatalf("ParseWorkflowsDirectory failed: %v", err)
		}

		if len(docs) != 3 {
			t.Errorf("Expected 3 workflow docs, got %d", len(docs))
		}

		// Check that we got all expected workflows
		names := make(map[string]bool)
		for _, doc := range docs {
			names[doc.Name] = true
		}

		expectedNames := []string{"CI", "Test", "Deploy"}
		for _, name := range expectedNames {
			if !names[name] {
				t.Errorf("Expected workflow '%s' not found", name)
			}
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		emptyDir := filepath.Join(tempDir, "empty")
		if err := os.MkdirAll(emptyDir, 0755); err != nil {
			t.Fatalf("Failed to create empty directory: %v", err)
		}

		docs, err := ParseWorkflowsDirectory(emptyDir)
		if err != nil {
			t.Fatalf("ParseWorkflowsDirectory failed: %v", err)
		}

		if len(docs) != 0 {
			t.Errorf("Expected 0 workflow docs for empty directory, got %d", len(docs))
		}
	})

	t.Run("directory with non-workflow files", func(t *testing.T) {
		mixedDir := filepath.Join(tempDir, "mixed")
		if err := os.MkdirAll(mixedDir, 0755); err != nil {
			t.Fatalf("Failed to create mixed directory: %v", err)
		}

		// Create workflow and non-workflow files
		workflows := map[string]string{
			"workflow.yml": `# @workflow.name: Valid
name: Valid`,
			"readme.md":  "# README",
			"script.sh":  "#!/bin/bash",
			"data.json":  "{}",
		}

		for name, content := range workflows {
			filePath := filepath.Join(mixedDir, name)
			if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
				t.Fatalf("Failed to create file %s: %v", name, err)
			}
		}

		docs, err := ParseWorkflowsDirectory(mixedDir)
		if err != nil {
			t.Fatalf("ParseWorkflowsDirectory failed: %v", err)
		}

		// Should only parse .yml/.yaml files
		if len(docs) != 1 {
			t.Errorf("Expected 1 workflow doc, got %d", len(docs))
		}

		if len(docs) > 0 && docs[0].Name != "Valid" {
			t.Errorf("Expected workflow name 'Valid', got '%s'", docs[0].Name)
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		_, err := ParseWorkflowsDirectory("/nonexistent/directory")
		// Should not return error for glob on non-existent directory
		// but should return empty slice
		if err != nil {
			t.Errorf("Expected no error for non-existent directory, got: %v", err)
		}
	})
}

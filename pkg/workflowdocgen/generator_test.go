package workflowdocgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMarkdownTable(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("generate table with basic workflow docs", func(t *testing.T) {
		docs := []*WorkflowDoc{
			{
				Name:        "CI Pipeline",
				Description: "Run tests and builds",
				Owners:      "team-platform",
				Tags:        "ci, automation",
				FileName:    "ci.yml",
			},
			{
				Name:        "Deploy",
				Description: "Deploy to production",
				Owners:      "team-ops",
				Tags:        "deployment",
				FileName:    "deploy.yml",
			},
		}

		outputPath := filepath.Join(tempDir, "output1.md")
		err := GenerateMarkdownTable(docs, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// Check header
		if !strings.Contains(output, "# Workflow Documentation") {
			t.Error("Expected header '# Workflow Documentation' not found")
		}

		// Check table header
		if !strings.Contains(output, "| Workflow | Description | Owners | Tags | File |") {
			t.Error("Expected table header not found")
		}

		// Check workflow entries
		if !strings.Contains(output, "CI Pipeline") {
			t.Error("Expected 'CI Pipeline' in output")
		}
		if !strings.Contains(output, "Run tests and builds") {
			t.Error("Expected 'Run tests and builds' in output")
		}
		if !strings.Contains(output, "team-platform") {
			t.Error("Expected 'team-platform' in output")
		}
		if !strings.Contains(output, "ci, automation") {
			t.Error("Expected 'ci, automation' in output")
		}

		if !strings.Contains(output, "Deploy") {
			t.Error("Expected 'Deploy' in output")
		}
		if !strings.Contains(output, "team-ops") {
			t.Error("Expected 'team-ops' in output")
		}
	})

	t.Run("generate table with missing fields", func(t *testing.T) {
		docs := []*WorkflowDoc{
			{
				Name:     "Test Workflow",
				FileName: "test.yml",
				// Other fields are empty
			},
		}

		outputPath := filepath.Join(tempDir, "output2.md")
		err := GenerateMarkdownTable(docs, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// Should use "-" for empty fields
		if !strings.Contains(output, "| Test Workflow | - | - | - | test.yml |") {
			t.Errorf("Expected row with '-' for empty fields, got:\n%s", output)
		}
	})

	t.Run("generate table with extended metadata", func(t *testing.T) {
		docs := []*WorkflowDoc{
			{
				Name:         "Advanced Workflow",
				Description:  "Complex workflow",
				FileName:     "advanced.yml",
				Params:       "environment, version",
				Results:      "artifact URLs",
				Permissions:  "write:packages",
				Requirements: "Node 18+",
			},
			{
				Name:        "Simple Workflow",
				Description: "No extra metadata",
				FileName:    "simple.yml",
			},
		}

		outputPath := filepath.Join(tempDir, "output3.md")
		err := GenerateMarkdownTable(docs, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// Check for detailed section
		if !strings.Contains(output, "## Detailed Workflow Information") {
			t.Error("Expected '## Detailed Workflow Information' section")
		}

		// Check for workflow name as heading
		if !strings.Contains(output, "### Advanced Workflow") {
			t.Error("Expected '### Advanced Workflow' heading")
		}

		// Check for extended metadata
		if !strings.Contains(output, "**Parameters:** environment, version") {
			t.Error("Expected parameters in output")
		}
		if !strings.Contains(output, "**Results:** artifact URLs") {
			t.Error("Expected results in output")
		}
		if !strings.Contains(output, "**Permissions:** write:packages") {
			t.Error("Expected permissions in output")
		}
		if !strings.Contains(output, "**Requirements:** Node 18+") {
			t.Error("Expected requirements in output")
		}

		// Simple workflow should not appear in detailed section
		if strings.Contains(output, "### Simple Workflow") {
			t.Error("Simple workflow should not appear in detailed section")
		}
	})

	t.Run("escape pipe characters in content", func(t *testing.T) {
		docs := []*WorkflowDoc{
			{
				Name:        "Test | Workflow",
				Description: "Uses | pipes",
				Owners:      "team | ops",
				Tags:        "tag1 | tag2",
				FileName:    "test.yml",
			},
		}

		outputPath := filepath.Join(tempDir, "output4.md")
		err := GenerateMarkdownTable(docs, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// Check that pipes are escaped
		if !strings.Contains(output, "Test \\| Workflow") {
			t.Error("Expected escaped pipe in name")
		}
		if !strings.Contains(output, "Uses \\| pipes") {
			t.Error("Expected escaped pipe in description")
		}
		if !strings.Contains(output, "team \\| ops") {
			t.Error("Expected escaped pipe in owners")
		}
		if !strings.Contains(output, "tag1 \\| tag2") {
			t.Error("Expected escaped pipe in tags")
		}
	})

	t.Run("empty workflow list", func(t *testing.T) {
		docs := []*WorkflowDoc{}

		outputPath := filepath.Join(tempDir, "output5.md")
		err := GenerateMarkdownTable(docs, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// Should still generate header and table structure
		if !strings.Contains(output, "# Workflow Documentation") {
			t.Error("Expected header even with no workflows")
		}
		if !strings.Contains(output, "| Workflow | Description | Owners | Tags | File |") {
			t.Error("Expected table header even with no workflows")
		}
	})

	t.Run("nil workflow list", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "output6.md")
		err := GenerateMarkdownTable(nil, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// Should still generate valid markdown
		if !strings.Contains(output, "# Workflow Documentation") {
			t.Error("Expected header even with nil workflows")
		}
	})

	t.Run("workflow with only some extended metadata", func(t *testing.T) {
		docs := []*WorkflowDoc{
			{
				Name:        "Partial Metadata",
				Description: "Workflow with partial metadata",
				FileName:    "partial.yml",
				Params:      "input param",
				// Results, Permissions, Requirements are empty
			},
		}

		outputPath := filepath.Join(tempDir, "output7.md")
		err := GenerateMarkdownTable(docs, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// Should have detailed section
		if !strings.Contains(output, "## Detailed Workflow Information") {
			t.Error("Expected detailed section with partial metadata")
		}

		// Should have parameters
		if !strings.Contains(output, "**Parameters:** input param") {
			t.Error("Expected parameters in output")
		}

		// Should not have empty fields
		if strings.Contains(output, "**Results:**") {
			t.Error("Should not include empty results field")
		}
		if strings.Contains(output, "**Permissions:**") {
			t.Error("Should not include empty permissions field")
		}
		if strings.Contains(output, "**Requirements:**") {
			t.Error("Should not include empty requirements field")
		}
	})

	t.Run("file write error - invalid path", func(t *testing.T) {
		docs := []*WorkflowDoc{
			{
				Name:     "Test",
				FileName: "test.yml",
			},
		}

		// Try to write to a path that doesn't exist and can't be created
		outputPath := "/invalid/nonexistent/path/output.md"
		err := GenerateMarkdownTable(docs, outputPath)
		if err == nil {
			t.Error("Expected error when writing to invalid path, got nil")
		}
	})

	t.Run("workflow name fallback to filename", func(t *testing.T) {
		docs := []*WorkflowDoc{
			{
				Name:         "", // Empty name
				Description:  "Test",
				FileName:     "unnamed.yml",
				Params:       "test param",
				Results:      "test result",
				Permissions:  "read",
				Requirements: "none",
			},
		}

		outputPath := filepath.Join(tempDir, "output8.md")
		err := GenerateMarkdownTable(docs, outputPath)
		if err != nil {
			t.Fatalf("GenerateMarkdownTable failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		output := string(content)

		// In detailed section, should use filename when name is empty
		if !strings.Contains(output, "### unnamed.yml") {
			t.Error("Expected filename as heading when name is empty")
		}
	})
}

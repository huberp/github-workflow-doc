package workflowdocgen

import (
	"fmt"
	"os"
	"strings"
)

// GenerateMarkdownTable generates a markdown table from workflow documentation
func GenerateMarkdownTable(docs []*WorkflowDoc, outputPath string) error {
	var sb strings.Builder

	// Write the header
	sb.WriteString("# Workflow Documentation\n\n")
	sb.WriteString("This document provides an overview of all GitHub workflows in this repository.\n\n")

	// Write the table header
	sb.WriteString("| Workflow | Description | Owners | Tags | File |\n")
	sb.WriteString("|----------|-------------|--------|------|------|\n")

	// Write each workflow as a row
	for _, doc := range docs {
		name := doc.Name
		if name == "" {
			name = "-"
		}

		description := doc.Description
		if description == "" {
			description = "-"
		}

		owners := doc.Owners
		if owners == "" {
			owners = "-"
		}

		tags := doc.Tags
		if tags == "" {
			tags = "-"
		}

		file := doc.FileName

		// Escape special markdown characters in content
		name = escapeMarkdown(name)
		description = escapeMarkdown(description)
		owners = escapeMarkdown(owners)
		tags = escapeMarkdown(tags)

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			name, description, owners, tags, file))
	}

	// Always add detailed section, but only show workflows with extended metadata
	sb.WriteString("\n## Detailed Workflow Information\n\n")

	hasAnyDetails := false
	for _, doc := range docs {
		if doc.Params == "" && doc.Results == "" && doc.Permissions == "" && doc.Requirements == "" {
			continue
		}

		hasAnyDetails = true
		workflowName := doc.Name
		if workflowName == "" {
			workflowName = doc.FileName
		}

		sb.WriteString(fmt.Sprintf("### %s\n\n", workflowName))

		if doc.Params != "" {
			sb.WriteString(fmt.Sprintf("**Parameters:** %s\n\n", doc.Params))
		}

		if doc.Results != "" {
			sb.WriteString(fmt.Sprintf("**Results:** %s\n\n", doc.Results))
		}

		if doc.Permissions != "" {
			sb.WriteString(fmt.Sprintf("**Permissions:** %s\n\n", doc.Permissions))
		}

		if doc.Requirements != "" {
			sb.WriteString(fmt.Sprintf("**Requirements:** %s\n\n", doc.Requirements))
		}
	}

	// If no workflows had extended metadata, add a note
	if !hasAnyDetails {
		sb.WriteString("_No workflows have extended metadata configured._\n\n")
	}

	// Write to file with readable permissions for collaborative environments
	// #nosec G306 - 0644 is intentional for collaborative environments
	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}

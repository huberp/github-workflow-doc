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
		
		// Escape pipe characters in content
		name = strings.ReplaceAll(name, "|", "\\|")
		description = strings.ReplaceAll(description, "|", "\\|")
		owners = strings.ReplaceAll(owners, "|", "\\|")
		tags = strings.ReplaceAll(tags, "|", "\\|")
		
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", 
			name, description, owners, tags, file))
	}

	// Add additional details section if workflows have extended metadata
	hasExtendedMeta := false
	for _, doc := range docs {
		if doc.Params != "" || doc.Results != "" || doc.Permissions != "" || doc.Requirements != "" {
			hasExtendedMeta = true
			break
		}
	}

	if hasExtendedMeta {
		sb.WriteString("\n## Detailed Workflow Information\n\n")
		
		for _, doc := range docs {
			if doc.Params == "" && doc.Results == "" && doc.Permissions == "" && doc.Requirements == "" {
				continue
			}
			
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
	}

	// Write to file
	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}

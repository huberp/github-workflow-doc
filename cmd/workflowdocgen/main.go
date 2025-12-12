package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/huberp/github-workflow-doc/pkg/workflowdocgen"
)

func main() {
	// Define flags
	workflowsDir := flag.String("workflows-dir", ".github/workflows", "Path to the workflows directory")
	outputFile := flag.String("output", "WORKFLOWS.md", "Path to the output markdown file")
	flag.Parse()

	// Check if workflows directory exists
	if _, err := os.Stat(*workflowsDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Workflows directory does not exist: %s\n", *workflowsDir)
		os.Exit(1)
	}

	// Parse all workflow files
	docs, err := workflowdocgen.ParseWorkflowsDirectory(*workflowsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing workflows: %v\n", err)
		os.Exit(1)
	}

	if len(docs) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: No workflow files found in %s\n", *workflowsDir)
	}

	// Generate markdown table
	absOutputPath, err := filepath.Abs(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output path: %v\n", err)
		os.Exit(1)
	}

	err = workflowdocgen.GenerateMarkdownTable(docs, absOutputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating markdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated workflow documentation at %s\n", absOutputPath)
	fmt.Printf("Documented %d workflow(s)\n", len(docs))
}

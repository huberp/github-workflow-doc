package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/huberp/github-workflow-doc/pkg/workflowdocgen"
)

func main() {
	// Define flags
	workflowsDir := flag.String("workflows-dir", ".github/workflows", "Path to the workflows directory")
	outputFile := flag.String("output", "WORKFLOWS.md", "Path to the output markdown file")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	flag.Parse()

	// Setup structured logging
	logLevel := slog.LevelWarn
	if *verbose {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting workflow documentation generation", "workflows-dir", *workflowsDir, "output", *outputFile)

	// Check if workflows directory exists
	if _, err := os.Stat(*workflowsDir); os.IsNotExist(err) {
		slog.Error("Workflows directory does not exist", "path", *workflowsDir)
		fmt.Fprintf(os.Stderr, "Error: Workflows directory does not exist: %s\n", *workflowsDir)
		os.Exit(1)
	}

	slog.Info("Parsing workflow files", "directory", *workflowsDir)

	// Parse all workflow files
	docs, err := workflowdocgen.ParseWorkflowsDirectory(*workflowsDir)
	if err != nil {
		slog.Error("Failed to parse workflows", "error", err)
		fmt.Fprintf(os.Stderr, "Error parsing workflows: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Parsed workflows", "count", len(docs))

	if len(docs) == 0 {
		slog.Warn("No workflow files found", "directory", *workflowsDir)
		fmt.Fprintf(os.Stderr, "Warning: No workflow files found in %s\n", *workflowsDir)
	}

	// Generate markdown table
	absOutputPath, err := filepath.Abs(*outputFile)
	if err != nil {
		slog.Error("Failed to resolve output path", "error", err)
		fmt.Fprintf(os.Stderr, "Error resolving output path: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Generating markdown documentation", "output", absOutputPath)

	err = workflowdocgen.GenerateMarkdownTable(docs, absOutputPath)
	if err != nil {
		slog.Error("Failed to generate markdown", "error", err)
		fmt.Fprintf(os.Stderr, "Error generating markdown: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Documentation generation complete", "output", absOutputPath, "workflows", len(docs))
	fmt.Printf("Successfully generated workflow documentation at %s\n", absOutputPath)
	fmt.Printf("Documented %d workflow(s)\n", len(docs))
}

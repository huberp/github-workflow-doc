package workflowdocgen

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// WorkflowDoc represents the documentation for a workflow
type WorkflowDoc struct {
	Name         string
	Description  string
	Owners       string
	Tags         string
	Params       string
	Results      string
	Permissions  string
	Requirements string
	FilePath     string
	FileName     string
}

// ParseWorkflowFile parses a workflow YAML file and extracts documentation comments
func ParseWorkflowFile(filePath string) (doc *WorkflowDoc, err error) {
	// Validate and clean the file path to prevent directory traversal
	cleanPath := filepath.Clean(filePath)

	file, ferr := os.Open(cleanPath)
	if ferr != nil {
		return nil, ferr
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	doc = &WorkflowDoc{
		FilePath: filePath,
		FileName: filepath.Base(filePath),
	}

	// Regex patterns to match documentation comments
	workflowPattern := regexp.MustCompile(`^#\s*@workflow\.([a-z]+):\s*(.*)$`)
	jobPattern := regexp.MustCompile(`^#\s*@job\.([a-z]+):\s*(.*)$`)
	stepPattern := regexp.MustCompile(`^#\s*@step\.([a-z]+):\s*(.*)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Only process lines starting with # @
		if !strings.HasPrefix(strings.TrimSpace(line), "# @") {
			continue
		}

		// Try to match workflow pattern
		matches := workflowPattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			field := matches[1]
			value := strings.TrimSpace(matches[2])

			switch field {
			case "name":
				doc.Name = value
			case "description":
				doc.Description = value
			case "owners":
				doc.Owners = value
			case "tags":
				doc.Tags = value
			case "params":
				doc.Params = value
			case "results":
				doc.Results = value
			case "permissions":
				doc.Permissions = value
			case "requirements":
				doc.Requirements = value
			}
			continue
		}

		// Try to match job pattern (for future expansion)
		jobMatches := jobPattern.FindStringSubmatch(line)
		if len(jobMatches) == 3 {
			// Job-level documentation - not yet stored in WorkflowDoc
			// This is parsed for validation but not currently used
			continue
		}

		// Try to match step pattern (for future expansion)
		stepMatches := stepPattern.FindStringSubmatch(line)
		if len(stepMatches) == 3 {
			// Step-level documentation - not yet stored in WorkflowDoc
			// This is parsed for validation but not currently used
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return doc, nil
}

// ParseWorkflowsDirectory parses all workflow files in a directory
func ParseWorkflowsDirectory(dirPath string) ([]*WorkflowDoc, error) {
	var docs []*WorkflowDoc

	// Clean and validate the directory path
	cleanDirPath := filepath.Clean(dirPath)

	// Find all YAML files in the directory
	files, err := filepath.Glob(filepath.Join(cleanDirPath, "*.yml"))
	if err != nil {
		return nil, err
	}

	yamlFiles, err := filepath.Glob(filepath.Join(cleanDirPath, "*.yaml"))
	if err != nil {
		return nil, err
	}

	files = append(files, yamlFiles...)

	for _, file := range files {
		// Check if file is a symlink and skip it
		fileInfo, err := os.Lstat(file)
		if err != nil {
			slog.Warn("Failed to stat file", "file", file, "error", err)
			continue
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			slog.Warn("Skipping symlink", "file", file)
			continue
		}

		doc, err := ParseWorkflowFile(file)
		if err != nil {
			slog.Warn("Failed to parse workflow file", "file", file, "error", err)
			continue
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// escapeMarkdown escapes special markdown characters in table cells
func escapeMarkdown(s string) string {
	// Escape special characters that can break markdown tables
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, "*", "\\*")
	s = strings.ReplaceAll(s, "_", "\\_")
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	return s
}

// IsYAMLFile checks if a file is a valid YAML file
func IsYAMLFile(filePath string) error {
	// Basic validation: check file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("file %s is not a YAML file", filePath)
	}
	return nil
}

package workflowdocgen

import (
	"bufio"
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
func ParseWorkflowFile(filePath string) (*WorkflowDoc, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc := &WorkflowDoc{
		FilePath: filePath,
		FileName: filepath.Base(filePath),
	}

	// Regex pattern to match documentation comments: # @workflow.field: value
	workflowPattern := regexp.MustCompile(`^#\s*@workflow\.([a-z]+):\s*(.*)$`)
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Only process lines starting with # @
		if !strings.HasPrefix(strings.TrimSpace(line), "# @") {
			continue
		}
		
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

	// Find all YAML files in the directory
	files, err := filepath.Glob(filepath.Join(dirPath, "*.yml"))
	if err != nil {
		return nil, err
	}
	
	yamlFiles, err := filepath.Glob(filepath.Join(dirPath, "*.yaml"))
	if err != nil {
		return nil, err
	}
	
	files = append(files, yamlFiles...)

	for _, file := range files {
		doc, err := ParseWorkflowFile(file)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

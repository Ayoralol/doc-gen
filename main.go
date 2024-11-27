package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	config, err := loadConfig("docs-config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	for _, fileConfig := range config.Files {
		yamlFiles, err := readYAMLFiles(fileConfig.Path)
		if err != nil {
			fmt.Printf("Error reading YAML files: %v\n", err)
			continue
		}
		for _, yamlFile := range yamlFiles {
			parsed, err := parseYAML(yamlFile)
			if err != nil {
				fmt.Printf("Error parsing YAML file %s: %v\n", yamlFile, err)
				continue
			}
			mdContent := toMarkdown(parsed, fileConfig.Type, filepath.Base(yamlFile), config.Repo, fileConfig.Path)
			mdPath := createOutputPath(fileConfig.DocPath, yamlFile, config.Output.Individual)
			if err := writeMarkdown(mdPath, mdContent); err != nil {
				fmt.Printf("Error writing Markdown file %s: %v\n", mdPath, err)
			}
		}
	}
}

func readYAMLFiles(basePath string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".yaml" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func parseYAML(file string) ([]OrderedYAML, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var node yaml.Node
	err = yaml.Unmarshal(content, &node)
	if err != nil {
		return nil, err
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return decodeOrderedYAML(node.Content[0]), nil
	}
	return nil, fmt.Errorf("invalid YAML structure")
}

func decodeOrderedYAML(node *yaml.Node) []OrderedYAML {
	var result []OrderedYAML
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i].Value
			val := decodeValue(node.Content[i+1])
			result = append(result, OrderedYAML{Key: key, Value: val})
		}
	} else {
	}
	return result
}

func decodeValue(node *yaml.Node) interface{} {
	switch node.Kind {
	case yaml.MappingNode:
		return decodeOrderedYAML(node)
	case yaml.SequenceNode:
		var result []interface{}
		for _, item := range node.Content {
			result = append(result, decodeValue(item))
		}
		return result
	case yaml.ScalarNode:
		return node.Value
	default:
		return nil
	}
}

func createOutputPath(docPath, filePath, outputBase string) string {
	relativePath := filepath.Base(filePath)
	mdPath := filepath.Join(outputBase, docPath, relativePath)
	return mdPath[:len(mdPath)-len(filepath.Ext(mdPath))] + ".md"
}

func toMarkdown(data []OrderedYAML, docType string, fileName string, repo string, filePath string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### %s\n", docType))
	sb.WriteString(fmt.Sprintf("> [%s](%vtree/main/%s/%v)\n\n", fileName, repo, filePath, fileName))
	sb.WriteString(yamlToMarkdown(data, ""))
	return sb.String()
}

func writeMarkdown(filePath string, content string) error {
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directories for %s: %v", filePath, err)
	}
	return os.WriteFile(filePath, []byte(content), 0644)
}

func yamlToMarkdown(data []OrderedYAML, indent string) string {
	var sb strings.Builder

	for _, item := range data {
		if strings.HasPrefix(item.Key, "-") {
			sb.WriteString(fmt.Sprintf("%s:\n", item.Key))
		} else {
			sb.WriteString(fmt.Sprintf("%s- %s:\n", indent, item.Key))
		}
		sb.WriteString(yamlToMarkdownWithFormatting(item.Key, item.Value, indent+"    "))
	}

	return cleanExtraNewlines(sb.String())
}

func yamlToMarkdownWithFormatting(key string, data interface{}, indent string) string {
	var sb strings.Builder

	switch value := data.(type) {
	case []OrderedYAML:
		sb.WriteString(yamlToMarkdown(value, indent))
	case []interface{}:
		for _, item := range value {
			if str, ok := item.(string); ok {
				sb.WriteString(fmt.Sprintf("%s- ***%s***\n", indent, str))
			} else {
				sb.WriteString(fmt.Sprintf("%s\n", yamlToMarkdownWithFormatting("", item, indent)))
			}
		}
	case string:
		if strings.Contains(key, "query") {
			sb.WriteString(fmt.Sprintf("%s```sql\n%s%v\n%s```\n", indent, indent, value, indent))
		} else {
			sb.WriteString(fmt.Sprintf("%s***%s***\n", indent, value))
		}
	default:
		if list, ok := value.([]interface{}); ok && len(list) > 0 {
			for _, item := range list {
				if str, ok := item.(string); ok {
					sb.WriteString(fmt.Sprintf("%s- ***%s***\n", indent, str))
				} else {
					sb.WriteString(fmt.Sprintf("%s- %s\n", indent, yamlToMarkdownWithFormatting("", item, indent+"    ")))
				}
			}
		} else {
			sb.WriteString(fmt.Sprintf("%s%v\n", indent, value))
		}
	}

	return sb.String()
}

func cleanExtraNewlines(input string) string {
	lines := strings.Split(input, "\n")
	var cleaned []string
	blank := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !blank {
				cleaned = append(cleaned, "")
				blank = true
			}
		} else {
			cleaned = append(cleaned, line)
			blank = false
		}
	}

	return strings.Join(cleaned, "\n")
}

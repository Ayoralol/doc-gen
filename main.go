package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// parseConfig reads a YAML file and parses it into a Config struct
func parseConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// Get all .yaml files in a folder
func getYamlFilesFromFolder(folderPath string) ([]string, error) {
	var yamlFiles []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to traverse folder %s: %w", folderPath, err)
	}

	return yamlFiles, nil
}

// Generate the aggregated documentation file
func generateAggregatedDocs(config *Config) error {
	aggregatedPath := config.Output.Aggregated

	// Open the aggregated file for writing
	file, err := os.Create(aggregatedPath)
	if err != nil {
		return fmt.Errorf("failed to create aggregated file: %w", err)
	}
	defer file.Close()

	// Write the main document header
	_, err = file.WriteString("# Documentation\n\n## Prometheus\n\n## Prometheus Scrape Configs\n\n")
	if err != nil {
		return fmt.Errorf("failed to write to aggregated file: %w", err)
	}

	// Iterate over the structure in docs-config.yaml
	for _, section := range config.Structure {
		for _, filePath := range section.Files {
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Error accessing path %s: %v\n", filePath, err)
				continue
			}

			var yamlFiles []string
			if fileInfo.IsDir() {
				// Get all .yaml files in the folder
				yamlFiles, err = getYamlFilesFromFolder(filePath)
				if err != nil {
					fmt.Printf("Error reading folder %s: %v\n", filePath, err)
					continue
				}
			} else {
				// Treat the path as a single file
				yamlFiles = []string{filePath}
			}

			for _, yamlFile := range yamlFiles {
				docFilePath := filepath.Join(config.Output.Individual, filepath.Base(yamlFile)+"-docs.md")

				// Read the individual doc file
				content, err := os.ReadFile(docFilePath)
				if err != nil {
					fmt.Printf("Error reading individual doc file %s: %v\n", docFilePath, err)
					continue
				}

				// Append only the job details (skip headers)
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "##") {
						continue // Skip section headers
					}
					_, err = file.WriteString(line + "\n")
					if err != nil {
						fmt.Printf("Error writing to aggregated file: %v\n", err)
						continue
					}
				}
			}
		}
	}

	fmt.Printf("Aggregated documentation generated at %s\n", aggregatedPath)
	return nil
}

func main() {
	// Parse the configuration file
	config, err := parseConfig("docs-config.yaml")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Generate individual documentation files
	for _, file := range config.Files {
		if file.Type == "prometheus-scrape-config" {
			// Check if the path is a directory
			fileInfo, err := os.Stat(file.Path)
			if err != nil {
				fmt.Printf("Error accessing path %s: %v\n", file.Path, err)
				continue
			}

			var yamlFiles []string
			if fileInfo.IsDir() {
				// Get all .yaml files in the folder
				yamlFiles, err = getYamlFilesFromFolder(file.Path)
				if err != nil {
					fmt.Printf("Error reading folder %s: %v\n", file.Path, err)
					continue
				}
			} else {
				// Treat the path as a single file
				yamlFiles = []string{file.Path}
			}

			// Process each .yaml file
			for _, yamlFile := range yamlFiles {
				scrapeConfig, err := parsePrometheusFile(yamlFile)
				if err != nil {
					fmt.Printf("Error parsing file %s: %v\n", yamlFile, err)
					continue
				}

				markdown, err := renderPrometheusMarkdown(scrapeConfig, yamlFile)
				if err != nil {
					fmt.Printf("Error rendering Markdown for file %s: %v\n", yamlFile, err)
					continue
				}

				err = writeMarkdownToFile(markdown, config.Output.Individual, filepath.Base(yamlFile))
				if err != nil {
					fmt.Printf("Error writing Markdown for file %s: %v\n", yamlFile, err)
					continue
				}

				fmt.Printf("Generated documentation for %s\n", yamlFile)
			}
		}
	}

	// Generate the aggregated documentation file
	err = generateAggregatedDocs(config)
	if err != nil {
		fmt.Println("Error generating aggregated docs:", err)
		return
	}
}

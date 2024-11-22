package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"text/template"
)

// PrometheusScrapeConfig represents a Prometheus scrape config file
type PrometheusScrapeConfig struct {
	// Flat structure fields
	JobName        string `yaml:"job_name"`
	ScrapeInterval string `yaml:"scrape_interval"`
	StaticConfigs  []struct {
		Targets []string `yaml:"targets"`
	} `yaml:"static_configs"`

	// Nested structure
	ScrapeConfigs []struct {
		JobName        string `yaml:"job_name"`
		ScrapeInterval string `yaml:"scrape_interval"`
		StaticConfigs  []struct {
			Targets []string `yaml:"targets"`
		} `yaml:"static_configs"`
	} `yaml:"scrape_configs"`
}

// Parse a Prometheus scrape config YAML file
func parsePrometheusFile(filePath string) (*PrometheusScrapeConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var scrapeConfig PrometheusScrapeConfig
	err = yaml.Unmarshal(data, &scrapeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Prometheus scrape config: %w", err)
	}

	return &scrapeConfig, nil
}

// Render Markdown for a Prometheus scrape config
func renderPrometheusMarkdown(config *PrometheusScrapeConfig, filePath string) (string, error) {
	// Load the template file
	tmplPath := "templates/prometheus.tmpl"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString("## Prometheus Scrape Configs\n\n")

	// Handle nested scrape_configs
	if len(config.ScrapeConfigs) > 0 {
		for _, job := range config.ScrapeConfigs {
			data := struct {
				FilePath       string
				JobName        string
				ScrapeInterval string
				Targets        []string
			}{
				FilePath:       filePath,
				JobName:        job.JobName,
				ScrapeInterval: job.ScrapeInterval,
				Targets:        job.StaticConfigs[0].Targets,
			}

			err = tmpl.Execute(&buf, data)
			if err != nil {
				return "", fmt.Errorf("failed to render template: %w", err)
			}
			buf.WriteString("\n")
		}
	} else {
		// Handle flat structure
		data := struct {
			FilePath       string
			JobName        string
			ScrapeInterval string
			Targets        []string
		}{
			FilePath:       filePath,
			JobName:        config.JobName,
			ScrapeInterval: config.ScrapeInterval,
			Targets:        config.StaticConfigs[0].Targets,
		}

		err = tmpl.Execute(&buf, data)
		if err != nil {
			return "", fmt.Errorf("failed to render template: %w", err)
		}
	}

	return buf.String(), nil
}

// Write Markdown content to a file
func writeMarkdownToFile(content, outputDir, fileName string) error {
	outputPath := filepath.Join(outputDir, fileName+"-docs.md")
	err := os.MkdirAll(outputDir, 0755) // Ensure the output directory exists
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	err = os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

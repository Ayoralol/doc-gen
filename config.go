package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type FileConfig struct {
	Path    string `yaml:"path"`
	DocPath string `yaml:"doc-path"`
	Type    string `yaml:"type"`
}

type OutputConfig struct {
	Individual string `yaml:"individual"`
	Aggregated string `yaml:"aggregated"`
}

type DocsConfig struct {
	Files  []FileConfig `yaml:"files"`
	Output OutputConfig `yaml:"output"`
	Repo   string       `yaml:"repository"`
}

func loadConfig(filePath string) (*DocsConfig, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config DocsConfig
	err = yaml.Unmarshal(content, &config)
	return &config, err
}

type OrderedYAML struct {
	Key   string
	Value interface{}
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkFileNames(dir string) error {
	var filesWithSpaces []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) == ".yaml" && strings.Contains(info.Name(), " ") {
			filesWithSpaces = append(filesWithSpaces, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(filesWithSpaces) > 0 {
		return fmt.Errorf("YAML files with spaces found: %v", filesWithSpaces)
	}

	return nil
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type MediaDirectory map[string][]string

var SupportedFileExtensions = []string{".mp3"}

func ReadFilesInDirectory(path string) (MediaDirectory, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	directory := make(MediaDirectory)

	for _, entry := range entries {
		fullFileName := fmt.Sprintf("%s/%s", path, entry.Name())
		if entry.IsDir() {
			subDirectory, e := ReadFilesInDirectory(fullFileName)
			if e != nil {
				return nil, e
			}

			for subPath, files := range subDirectory {
				directory[subPath] = append(directory[subPath], files...)
			}

			continue
		}
		if slices.Contains(SupportedFileExtensions, filepath.Ext(fullFileName)) {
			directory[path] = append(directory[path], entry.Name())
		}
	}

	return directory, nil
}

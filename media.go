package main

import (
	"fmt"
	"os"
)

type MediaFile struct {
	name string
	path string
}

func ReadFilesInDirectory(dir string) ([]MediaFile, error) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	files := make([]MediaFile, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			f, e := ReadFilesInDirectory(fmt.Sprintf("%s/%s", dir, entry.Name()))
			if e != nil {
				return nil, e
			}

			files = append(files, f...)
			continue
		}
		files = append(files, MediaFile{name: entry.Name(), path: dir})
	}

	return files, nil
}

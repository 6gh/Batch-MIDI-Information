package main

import (
	"io/fs"
	"path/filepath"
)

func findMIDIs(dir string) ([]string, error) {
	midis := []string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".mid" {
			println("Found MIDI:", path)
			midis = append(midis, path)
		}

		return nil
	})
	return midis, err
}

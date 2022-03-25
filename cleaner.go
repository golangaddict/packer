package packer

import (
	"os"
	"path/filepath"
)

type Cleaner []string

func (c Cleaner) Run(path string) error {
	for _, cleanPath := range c {
		entries, err := os.ReadDir(cleanPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if err := os.RemoveAll(filepath.Join(cleanPath, entry.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

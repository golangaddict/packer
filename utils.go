package packer

import (
	"io/fs"
	"os"
	"path/filepath"
)

func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return stat.IsDir(), nil
}

func getFilesFromPath(path string) (s []string, err error) {
	isDir, err := isDir(path)
	if err != nil {
		return nil, err
	}

	if !isDir {
		return []string{path}, nil
	}

	return s, filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		s = append(s, path)

		return nil
	})
}

func sliceContains(s []string, v string) bool {
	for _, x := range s {
		if v == x {
			return true
		}
	}

	return false
}

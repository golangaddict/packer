package packer

import (
	"io/fs"
	"os"
	"path/filepath"
)

func IsDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return stat.IsDir(), nil
}

func GetFilesFromDir(path string) (s []string, err error) {
	return s, filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		s = append(s, path)

		return nil
	})
}

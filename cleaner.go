package packer

import "os"

type Cleaner []string

func (c Cleaner) Run(path string) error {
	for _, cleanPath := range c {
		if err := os.RemoveAll(cleanPath); err != nil {
			return err
		}
	}

	return nil
}

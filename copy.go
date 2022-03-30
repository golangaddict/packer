package packer

import "os/exec"

type Copier []struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (c Copier) Run(path string) error {
	for _, v := range c {
		if err := exec.Command("cp", "-r", v.From, v.To).Run(); err != nil {
			return err
		}
	}

	return nil
}

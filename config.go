package main

type Config []struct {
	Name string   `json:"name" validate:"name"`
	JS   JSConfig `json:"js"`
}

type JSConfig struct {
	Libs        []string `json:"libs"`
	Deps        []string `json:"deps"`
	ModulePaths []string `json:"module_paths"`
	Modules     []string `json:"modules"`
	Output      string   `json:"output"`
}

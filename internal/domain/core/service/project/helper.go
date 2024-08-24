package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func (*Project) createProjectDir(dirParam string) (string, error) {
	projectPath, err := filepath.Abs(dirParam)
	if err != nil {
		return "", fmt.Errorf("filepath: abs: %w", err)
	}

	stat, err := os.Stat(projectPath)
	if !os.IsNotExist(err) {
		if !stat.IsDir() {
			return "", fmt.Errorf("stat: is dir: %w", errors.New("dir must be folder"))
		}

		projectFiles := filepath.Join(projectPath, "*")

		glob, err := filepath.Glob(projectFiles)
		if err != nil {
			return "", fmt.Errorf("filepath: glob: %w", err)
		}

		if len(glob) > 0 {
			return "", fmt.Errorf("check project folder is empty: %w", errors.New("project folder must be empty"))
		}
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(projectPath, 0o755)
		if err != nil {
			return "", fmt.Errorf("os: mkdir all: %w", err)
		}
	}

	err = os.Chdir(projectPath)
	if err != nil {
		return "", fmt.Errorf("os: chdir: %w", err)
	}

	return projectPath, nil
}

func (*Project) createHexagoConfigs(projectPath string) error {
	hexagoDir := filepath.Join(projectPath, ".hexago")

	err := os.MkdirAll(hexagoDir, 0o755)
	if err != nil {
		return fmt.Errorf("os: mkdir all: %w", err)
	}

	configPath := filepath.Join(hexagoDir, "config.yaml")

	configContent, err := assets.ReadFile("assets/config.yaml")
	if err != nil {
		return fmt.Errorf("assets: read file: %w", err)
	}

	err = os.WriteFile(configPath, configContent, 0o644)
	if err != nil {
		return fmt.Errorf("os: write file: %w", err)
	}

	templatesPath := filepath.Join(hexagoDir, "templates")
	err = os.MkdirAll(templatesPath, 0o755)
	if err != nil {
		return fmt.Errorf("os: mkdir all: %w", err)
	}

	return nil
}

func (*Project) createProjectSubDirs() error {
	dirs := []string{
		"cmd",
		filepath.Join("internal", "domain", "core", "application"),
		filepath.Join("internal", "domain", "core", "dto"),
		filepath.Join("internal", "domain", "core", "model"),
		filepath.Join("internal", "domain", "core", "port"),
		filepath.Join("internal", "domain", "core", "service"),
		filepath.Join("internal", "infrastructure"),
		filepath.Join("internal", "pkg"),
		"pkg",
		"config",
		"schemas",
		"scripts",
		"doc",
	}

	for i := range dirs {
		err := os.MkdirAll(dirs[i], 0o755)
		if err != nil {
			return fmt.Errorf("os: mkdir all: %w", err)
		}
	}

	return nil
}

var instanceNameRgx = regexp.MustCompile(`^[A-Z][A-Za-z0-9]{0,}$`)

func (*Project) validateInstanceName(instanceType, instanceName string) error {
	if !instanceNameRgx.MatchString(instanceName) {
		return fmt.Errorf("invalid %[1]s name: %[2]s, %[1]s name must be \"PascalCase\"", instanceType, instanceName)
	}
	return nil
}

var pkgNameRgx = regexp.MustCompile(`^[a-z][a-z0-9]{0,}$`)

func (*Project) validatePkgName(pkgName string) error {
	if !pkgNameRgx.MatchString(pkgName) {
		return fmt.Errorf("invalid package name: %s, package name must be \"lowercase\"", pkgName)
	}
	return nil
}

var pkgCmdRgx = regexp.MustCompile(`^[a-z][a-z0-9]{0,}$`)

func (*Project) validateEntryPointName(entryPointName string) error {
	if !pkgNameRgx.MatchString(entryPointName) {
		return fmt.Errorf("invalid entry point name: %s, entry point name must be \"lowercase\" \"lower_case\" or \"lower-case\"", entryPointName)
	}
	return nil
}

package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
)

func (*Project) getTerminalName() (string, error) {
	if runtime.GOOS == "windows" {
		return "powershell", nil
	}
	shellPath, ok := os.LookupEnv("SHELL")
	if !ok {
		return "", fmt.Errorf("SHELL env not found")
	}
	return strings.TrimSpace(shellPath), nil
}

func (*Project) createProjectDir(dirParam string) (string, error) {
	projectPath, err := filepath.Abs(dirParam)
	if err != nil {
		return "", fmt.Errorf("filepath: abs: %w", err)
	}

	stat, err := os.Stat(projectPath)
	if !os.IsNotExist(err) && !stat.IsDir() {
		return "", fmt.Errorf("stat: is dir: %w", dto.ErrDirMustBeFolder)
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

func (*Project) addGitignore(projectPath string) error {
	configContent, err := assets.ReadFile("assets/.gitignore")
	if err != nil {
		return fmt.Errorf("assets: read file: %w", err)
	}

	err = os.WriteFile(filepath.Join(projectPath, ".gitignore"), configContent, 0o644)
	if err != nil {
		return fmt.Errorf("os: write file: %w", err)
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

func (*Project) ValidateInstanceName(instanceName string) error {
	if !instanceNameRgx.MatchString(instanceName) {
		return dto.ErrInvalidInstanceName
	}
	return nil
}

var pkgNameRgx = regexp.MustCompile(`^[a-z][a-z0-9]{0,}$`)

func (*Project) ValidatePkgName(pkgName string) error {
	if !pkgNameRgx.MatchString(pkgName) {
		return dto.ErrInvalidPkgName
	}
	return nil
}

var pkgCmdRgx = regexp.MustCompile(`^[a-z][a-z0-9\-]{0,}$`)

func (*Project) ValidateEntryPointName(entryPointName string) error {
	if !pkgNameRgx.MatchString(entryPointName) {
		return dto.ErrInvalidCmdName
	}
	return nil
}

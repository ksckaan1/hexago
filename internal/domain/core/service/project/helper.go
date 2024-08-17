package project

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
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

	err = os.WriteFile(configPath, nil, 0o644)
	if err != nil {
		return fmt.Errorf("os: write file: %w", err)
	}

	return nil
}

func (*Project) initGoModule(ctx context.Context, moduleName string) error {
	cmd := exec.CommandContext(ctx, "go", "mod", "init", moduleName)
	stdErr := &bytes.Buffer{}
	cmd.Stderr = stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd: run: %w", errors.New(strings.TrimSpace(stdErr.String())))
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
		filepath.Join("internal", "infrastructure", "repository"),
		filepath.Join("internal", "util"),
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

func (p *Project) isDomainExist(ctx context.Context, targetDomain string) error {
	domains, err := p.GetAllDomains(ctx)
	if err != nil {
		return fmt.Errorf("get all domains: %w", err)
	}

	if !slices.Contains(domains, targetDomain) {
		return fmt.Errorf("target domain not found: %s", targetDomain)
	}

	return nil
}

func (p *Project) isServiceExist(ctx context.Context, targetDomain, targetService string) error {
	services, err := p.GetAllServices(ctx, targetDomain)
	if err != nil {
		return fmt.Errorf("get all services: %w", err)
	}

	if !slices.Contains(services, targetService) {
		return fmt.Errorf("target service not found: %s", targetService)
	}

	return nil
}

var serviceNameRgx = regexp.MustCompile(`^[A-Z][A-Za-z0-9]{0,}$`)

func (*Project) validateServiceName(serviceName string) error {
	if !serviceNameRgx.MatchString(serviceName) {
		return fmt.Errorf("invalid service name: %s, service name must be PascalCase", serviceName)
	}
	return nil
}

var pkgNameRgx = regexp.MustCompile(`^[a-z][a-z0-9]{0,}$`)

func (*Project) validatePkgName(pkgName string) error {
	if !pkgNameRgx.MatchString(pkgName) {
		return fmt.Errorf("invalid package name: %s, package name must be lowercase", pkgName)
	}
	return nil
}

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
	"text/template"
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

	// create default templates
	templates := []string{
		"std_application.tmpl",
		"std_service.tmpl",
		"do_application.tmpl",
		"do_service.tmpl",
	}

	for i := range templates {
		tmpl, err2 := assets.ReadFile(fmt.Sprintf("assets/%s", templates[i]))
		if err2 != nil {
			return fmt.Errorf("assets: read file: %w", err2)
		}

		tmplPath := filepath.Join(templatesPath, templates[i])
		err = os.WriteFile(tmplPath, tmpl, 0o644)
		if err != nil {
			return fmt.Errorf("os: write file: %w", err)
		}
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

func (p *Project) generateServiceFile(servicePath, serviceName, pkgName string) (string, error) {
	serviceFile := filepath.Join(servicePath, fmt.Sprintf("%s.go", pkgName))

	serviceTemplatePath := filepath.Join(".hexago", "templates", fmt.Sprintf("%s_service.tmpl", p.cfg.GetServiceTemplate()))

	serviceTemplate, err := template.ParseFiles(serviceTemplatePath)
	if err != nil {
		return "", fmt.Errorf("template: parse files: %w", err)
	}

	buf := &bytes.Buffer{}
	err = serviceTemplate.Execute(buf, map[string]any{
		"ServiceName": serviceName,
		"PkgName":     pkgName,
	})
	if err != nil {
		return "", fmt.Errorf("template: execute: %w", err)
	}

	err = os.WriteFile(serviceFile, buf.Bytes(), 0o644)
	if err != nil {
		return "", fmt.Errorf("os: write file: %w", err)
	}

	return serviceFile, nil
}

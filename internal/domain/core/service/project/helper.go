package project

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/samber/lo"
	"golang.org/x/mod/modfile"
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

func (p *Project) generateServiceFile(ctx context.Context, targetDomain, servicePath, serviceName, pkgName, portParam string) (string, error) {
	serviceFile := filepath.Join(servicePath, fmt.Sprintf("%s.go", pkgName))

	portName, portDomain, portPath, implementation, err := p.generateImplementation(ctx, targetDomain, serviceName, portParam)
	if err != nil {
		return "", fmt.Errorf("generate implementation: %w", err)
	}

	serviceTemplate, err := p.parseTemplate(p.cfg.GetServiceTemplate(), "service")
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	buf := &bytes.Buffer{}
	err = serviceTemplate.Execute(buf, map[string]any{
		"ServiceName":        serviceName,
		"PkgName":            pkgName,
		"PortImplementation": implementation,
		"PortDomain":         portDomain,
		"TargetDomain":       targetDomain,
		"PortPath":           portPath,
		"PortName":           portName,
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

func (p *Project) parseTemplate(templateMode, templateName string) (*template.Template, error) {
	switch templateMode {
	case "std", "do":
		f, err := assets.ReadFile(fmt.Sprintf("assets/templates/%s_%s.tmpl", templateMode, templateName))
		if err != nil {
			return nil, fmt.Errorf("assets: read file: %w", err)
		}

		tmpl, err := template.New("").Parse(string(f))
		if err != nil {
			return nil, fmt.Errorf("template: parse: %w", err)
		}

		return tmpl, nil
	default:
		templatePath := filepath.Join(".hexago", "templates", fmt.Sprintf("%s_%s.tmpl", templateMode, templateName))

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return nil, fmt.Errorf("template: parse files: %w", err)
		}

		return tmpl, nil
	}

}

var rgxInterfaces = regexp.MustCompile(`(?m)^type\s([A-Z][a-zA-Z0-9]+)\sinterface`)

func (p *Project) parseInterfaces(interfaceFile string) ([]string, error) {
	f, err := os.Open(interfaceFile)
	if err != nil {
		return nil, fmt.Errorf("os: open: %w", err)
	}
	defer f.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, fmt.Errorf("io: copy: %w", err)
	}

	submatches := rgxInterfaces.FindAllSubmatch(buf.Bytes(), -1)

	interfaces := lo.Map(submatches, func(m [][]byte, _ int) string {
		return string(m[1])
	})

	return interfaces, nil
}

type PortValue struct {
	Name   string
	Domain string
}

var rgxDomainPort = regexp.MustCompile(`^([a-z][a-z0-9]{0,}):([A-Z][\w]{0,})$`)

var rgxPort = regexp.MustCompile(`^([A-Z][\w]{0,})$`)

func (*Project) getPort(targetDomain, portName string) (*PortValue, error) {
	if portName == "" {
		return &PortValue{}, nil
	}

	for _, sm := range rgxDomainPort.FindAllStringSubmatch(portName, -1) {
		return &PortValue{
			Name:   sm[2],
			Domain: sm[1],
		}, nil
	}

	if !rgxPort.MatchString(portName) {
		return nil, fmt.Errorf("invalid port name: %s", portName)
	}

	return &PortValue{
		Name:   portName,
		Domain: targetDomain,
	}, nil
}

func (p *Project) generateImplementation(ctx context.Context, targetDomain, serviceName, portParam string) (portName, portDomain, portPath, implementation string, err error) {
	portValue, err := p.getPort(targetDomain, portParam)
	if err != nil {
		return "", "", "", "", fmt.Errorf("get port: %w", err)
	}

	if portValue.Name == "" {
		return "", "", "", "", nil
	}

	f, err := os.Open("go.mod")
	if err != nil {
		return "", "", "", "", fmt.Errorf("module file not found: go.mod")
	}
	defer f.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, f)
	if err != nil {
		return "", "", "", "", fmt.Errorf("io: copy: %w", err)
	}

	modFile, err := modfile.Parse("go.mod", buf.Bytes(), nil)
	if err != nil {
		return "", "", "", "", fmt.Errorf("modfile parse: %w", err)
	}

	portPath = filepath.Join(modFile.Module.Mod.Path, "internal", "domain", portValue.Domain, "port")

	cmd := exec.CommandContext(ctx,
		"impl",
		fmt.Sprintf("%s *%s", strings.ToLower(string(serviceName[0])), serviceName),
		fmt.Sprintf("%s.%s", portPath, portValue.Name),
	)
	stdOut, stdErr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdOut, stdErr

	err = cmd.Run()
	if err != nil {
		return "", "", "", "", fmt.Errorf("impl: %s", stdErr.String())
	}

	return portValue.Name, portValue.Domain, portPath, stdOut.String(), nil
}

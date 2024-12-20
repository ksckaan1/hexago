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
	"strings"
	"text/template"

	"github.com/samber/lo"

	"github.com/ksckaan1/hexago/internal/customerrors"
)

type TemplateType string

const (
	TTService     TemplateType = "service"
	TTApplication TemplateType = "application"
	TTPackage     TemplateType = "package"
	TTInfra       TemplateType = "infra"
)

func (p *Project) generateGoInitFile(ctx context.Context, dir, structName, pkgName, portParam string, tt TemplateType, assertInterface bool) (string, error) {
	targetFilePath := filepath.Join(dir, fmt.Sprintf("%s.go", pkgName))

	implementationDetails, err := p.generateImplementation(ctx, structName, portParam)
	if err != nil {
		return "", fmt.Errorf("generate implementation: %w", err)
	}

	var templateMode string

	switch tt {
	case TTService:
		templateMode = p.cfg.GetServiceTemplate()
	case TTApplication:
		templateMode = p.cfg.GetApplicationTemplate()
	case TTInfra:
		templateMode = p.cfg.GetInfrastructureTemplate()
	case TTPackage:
		templateMode = p.cfg.GetPackageTemplate()
	default:
		return "", fmt.Errorf("invalid template type: %s", string(tt))
	}

	packageTemplate, err := p.parseTemplate(templateMode, string(tt))
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var (
		implementation string
		importPath     string
		interfaceName  string
	)

	if implementationDetails != nil {
		implementation = implementationDetails.Implementation
		importPath = implementationDetails.ImportPath
		interfaceName = implementationDetails.InterfaceName
	}

	buf := &bytes.Buffer{}
	err = packageTemplate.Execute(buf, map[string]any{
		"StructName":      structName,
		"PkgName":         pkgName,
		"Implementation":  implementation,
		"ImportPath":      importPath,
		"InterfaceName":   interfaceName,
		"AssertInterface": assertInterface,
	})
	if err != nil {
		return "", fmt.Errorf("template: execute: %w", customerrors.ErrTemplateCanNotExecute{Message: err.Error()})
	}

	err = os.WriteFile(targetFilePath, buf.Bytes(), 0o600)
	if err != nil {
		return "", fmt.Errorf("os: write file: %w", err)
	}

	return targetFilePath, nil
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
			return nil, fmt.Errorf("template: parse files: %w", errors.Join(err, customerrors.ErrTemplateCanNotParsed))
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
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			fmt.Println(err2)
		}
	}()

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

type ImplementationDetail struct {
	InterfaceName  string
	ImportPath     string
	Implementation string
}

func (p *Project) generateImplementation(ctx context.Context, instanceName, interfaceParam string) (*ImplementationDetail, error) {
	if interfaceParam == "" {
		return nil, nil
	}

	interfaceInfo, err := p.getInterfaceInfo(ctx, interfaceParam)
	if err != nil {
		return nil, fmt.Errorf("get interface info: %w", err)
	}

	if interfaceInfo == nil {
		return nil, nil
	}

	cmd := exec.Command(
		"impl",
		fmt.Sprintf("%s *%s", strings.ToLower(string(instanceName[0])), instanceName),
		interfaceInfo.ImplementParam,
	)

	stdOut, stdErr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdOut, stdErr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("impl: %s", stdErr.String())
	}

	return &ImplementationDetail{
		InterfaceName:  interfaceInfo.InterfaceName,
		ImportPath:     interfaceInfo.ImportPath,
		Implementation: stdOut.String(),
	}, nil
}

var (
	rgxPortParam   = regexp.MustCompile(`^([A-Z]\w*)$`)
	rgxNormalParam = regexp.MustCompile(`^(\S+)\.([A-Z]\w*)$`)
)

type InterfaceInfo struct {
	InterfaceName  string
	ImplementParam string
	ImportPath     string
	IsInDomain     bool
}

func (p *Project) getInterfaceInfo(_ context.Context, interfaceParam string) (*InterfaceInfo, error) {
	isPortParam := rgxPortParam.MatchString(interfaceParam)
	isNormalParam := rgxNormalParam.MatchString(interfaceParam)

	if !(isPortParam || isNormalParam) {
		return nil, customerrors.ErrInvalidPortName{PortName: interfaceParam}
	}

	if isNormalParam {
		sm := rgxNormalParam.FindStringSubmatch(interfaceParam)
		return &InterfaceInfo{
			InterfaceName:  sm[2],
			ImplementParam: interfaceParam,
			ImportPath:     sm[1],
		}, nil
	}

	moduleName, err := p.GetModuleName()
	if err != nil {
		return nil, fmt.Errorf("get module name: %w", err)
	}

	interfacePath := filepath.Join(moduleName, "internal", "port")

	return &InterfaceInfo{
		InterfaceName:  interfaceParam,
		ImplementParam: fmt.Sprintf("%s.%s", interfacePath, interfaceParam),
		ImportPath:     interfacePath,
		IsInDomain:     true,
	}, nil
}

func (p *Project) formatGoFiles(ctx context.Context, goFilePaths ...string) error {
	if len(goFilePaths) == 0 {
		return fmt.Errorf("no go files given")
	}

	args := []string{
		"-w",
	}

	args = append(args, goFilePaths...)

	cmd := exec.CommandContext(ctx, "gofmt", args...)
	stdErr := &bytes.Buffer{}
	cmd.Stderr = stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd: run: %w", customerrors.ErrFormatGoFile{Message: stdErr.String()})
	}

	return nil
}

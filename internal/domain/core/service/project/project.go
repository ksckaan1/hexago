package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/samber/do"
	"github.com/samber/lo"
)

var _ port.ProjectService = (*Project)(nil)

type Project struct {
	cfg port.ConfigService
}

func New(i *do.Injector) (port.ProjectService, error) {
	return &Project{
		cfg: do.MustInvoke[port.ConfigService](i),
	}, nil
}

func (p *Project) InitNewProject(ctx context.Context, params dto.InitNewProjectParams) error {
	projectPath, err := p.createProjectDir(params.ProjectDirectory)
	if err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}

	err = p.createHexagoConfigs(projectPath)
	if err != nil {
		return fmt.Errorf("create hexago configs: %w", err)
	}

	err = p.initGoModule(ctx, params.ModuleName)
	if err != nil {
		return fmt.Errorf("init go module: %w", err)
	}

	err = p.createProjectSubDirs()
	if err != nil {
		return fmt.Errorf("create project dirs: %w", err)
	}

	return nil
}

func (p *Project) GetAllDomains(ctx context.Context) ([]string, error) {
	domainLocation := filepath.Join("internal", "domain")

	domainCandidatePaths, err := filepath.Glob(filepath.Join(domainLocation, "*"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	domainPaths := lo.Filter(domainCandidatePaths, func(d string, _ int) bool {
		stat, err2 := os.Stat(d)
		return err2 == nil && stat.IsDir()
	})

	domains := lo.Map(domainPaths, func(d string, _ int) string {
		return filepath.Base(d)
	})

	return domains, nil
}

func (p *Project) CreateDomain(ctx context.Context, targetDomain string) error {
	err := p.isDomainExist(ctx, targetDomain)
	if err == nil {
		return fmt.Errorf("domain already exist: %s", targetDomain)
	}

	domainPath := filepath.Join("internal", "domain", targetDomain)

	domainDirs := []string{
		filepath.Join(domainPath, "application"),
		filepath.Join(domainPath, "dto"),
		filepath.Join(domainPath, "model"),
		filepath.Join(domainPath, "port"),
		filepath.Join(domainPath, "service"),
	}

	for i := range domainDirs {
		err = os.MkdirAll(domainDirs[i], 0o755)
		if err != nil {
			return fmt.Errorf("os: mkdir all: %w", err)
		}
	}

	return nil
}

func (p *Project) GetAllApplications(ctx context.Context, targetDomain string) ([]string, error) {
	err := p.isDomainExist(ctx, targetDomain)
	if err != nil {
		return nil, fmt.Errorf("is domain exist: %w", err)
	}

	applicationsPath := filepath.Join("internal", "domain", targetDomain, "application")

	applicationCandidatePaths, err := filepath.Glob(filepath.Join(applicationsPath, "*"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	applicationPaths := lo.Filter(applicationCandidatePaths, func(d string, _ int) bool {
		stat, err2 := os.Stat(d)
		return err2 == nil && stat.IsDir()
	})

	applications := lo.Map(applicationPaths, func(d string, _ int) string {
		return filepath.Base(d)
	})

	return applications, nil
}

func (p *Project) GetAllEntryPoints(ctx context.Context) ([]string, error) {
	cmdLocation := filepath.Join("cmd")

	cmdCandidatePaths, err := filepath.Glob(filepath.Join(cmdLocation, "*"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	cmdPaths := lo.Filter(cmdCandidatePaths, func(d string, _ int) bool {
		stat, err2 := os.Stat(d)
		return err2 == nil && stat.IsDir()
	})

	cmds := lo.Map(cmdPaths, func(d string, _ int) string {
		return filepath.Base(d)
	})

	return cmds, nil
}

func (p *Project) GetAllPorts(ctx context.Context, targetDomain string) ([]string, error) {
	err := p.isDomainExist(ctx, targetDomain)
	if err != nil {
		return nil, fmt.Errorf("is domain exist: %w", err)
	}

	portsPath := filepath.Join("internal", "domain", targetDomain, "port")

	portCandidatePaths, err := filepath.Glob(filepath.Join(portsPath, "*"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	portPaths := lo.Filter(portCandidatePaths, func(d string, _ int) bool {
		stat, err2 := os.Stat(d)
		return err2 == nil && !stat.IsDir()
	})

	ports := lo.Map(portPaths, func(d string, _ int) string {
		return filepath.Base(d)
	})

	return ports, nil
}

func (p *Project) GetAllServices(ctx context.Context, targetDomain string) ([]string, error) {
	err := p.isDomainExist(ctx, targetDomain)
	if err != nil {
		return nil, fmt.Errorf("is domain exist: %w", err)
	}

	servicesPath := filepath.Join("internal", "domain", targetDomain, "service")

	serviceCandidatePaths, err := filepath.Glob(filepath.Join(servicesPath, "*"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	servicePaths := lo.Filter(serviceCandidatePaths, func(d string, _ int) bool {
		stat, err2 := os.Stat(d)
		return err2 == nil && stat.IsDir()
	})

	services := lo.Map(servicePaths, func(d string, _ int) string {
		return filepath.Base(d)
	})

	return services, nil
}

func (p *Project) CreateService(ctx context.Context, targetDomain, serviceName, pkgName string) (string, error) {
	err := p.isDomainExist(ctx, targetDomain)
	if err != nil {
		return "", fmt.Errorf("is domain exist: %w", err)
	}

	err = p.validateServiceName(serviceName)
	if err != nil {
		return "", fmt.Errorf("validate service name: %w", err)
	}

	if pkgName == "" {
		pkgName = strings.ToLower(serviceName)
	}

	err = p.validatePkgName(pkgName)
	if err != nil {
		return "", fmt.Errorf("validate pkg name: %w", err)
	}

	err = p.isServiceExist(ctx, targetDomain, pkgName)
	if err == nil {
		return "", fmt.Errorf("service already exist: %s in %s", pkgName, targetDomain)
	}

	servicePath := filepath.Join("internal", "domain", targetDomain, "service", pkgName)

	err = os.MkdirAll(servicePath, 0o755)
	if err != nil {
		return "", fmt.Errorf("os: mkdir all: %w", err)
	}

	serviceFile, err := p.generateServiceFile(servicePath, serviceName, pkgName)
	if err != nil {
		return "", fmt.Errorf("generate service file: %w", err)
	}
	// TODO: implement a port for the generated service

	return serviceFile, nil
}

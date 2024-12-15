package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/samber/lo"

	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

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

func (p *Project) CreateService(ctx context.Context, params model.CreateServiceParams) (string, error) {
	err := p.ValidateInstanceName(params.StructName)
	if err != nil {
		return "", fmt.Errorf("validate instance name: %w", err)
	}

	if params.PackageName == "" {
		params.PackageName = strings.ToLower(params.StructName)
	}

	err = p.ValidatePkgName(params.PackageName)
	if err != nil {
		return "", fmt.Errorf("validate pkg name: %w", err)
	}

	err = p.isDomainExist(ctx, params.TargetDomain)
	if err != nil {
		return "", fmt.Errorf("is domain exist: %w", err)
	}

	err = p.isServiceExist(ctx, params.TargetDomain, params.PackageName)
	if err == nil {
		return "", fmt.Errorf("is service exist: %w", customerrors.ErrAlreadyExist)
	}

	serviceDir := filepath.Join("internal", "domain", params.TargetDomain, "service", params.PackageName)

	err = os.MkdirAll(serviceDir, 0o755)
	if err != nil {
		return "", fmt.Errorf("os: mkdir all: %w", err)
	}

	serviceFile, err := p.generateGoInitFile(
		ctx,
		serviceDir,
		params.StructName,
		params.PackageName,
		params.PortParam,
		TTService,
		params.AssertInterface,
	)
	if err != nil {
		err2 := os.RemoveAll(serviceDir)
		if err2 != nil {
			return "", fmt.Errorf("os: remove all: %w", err2)
		}
		return "", fmt.Errorf("generate service file: %w", err)
	}

	err = p.formatGoFiles(ctx, serviceFile)
	if err != nil {
		return "", fmt.Errorf("format go files: %w", err)
	}

	return serviceFile, nil
}

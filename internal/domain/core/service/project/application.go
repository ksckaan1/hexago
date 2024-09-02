package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/samber/lo"
)

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

func (p *Project) isApplicationExist(ctx context.Context, targetDomain, targetApplication string) error {
	applications, err := p.GetAllApplications(ctx, targetDomain)
	if err != nil {
		return fmt.Errorf("get all applications: %w", err)
	}

	if !slices.Contains(applications, targetApplication) {
		return fmt.Errorf("target application not found: %s", targetApplication)
	}

	return nil
}

func (p *Project) CreateApplication(ctx context.Context, params dto.CreateApplicationParams) (string, error) {
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

	err = p.isApplicationExist(ctx, params.TargetDomain, params.PackageName)
	if err == nil {
		return "", fmt.Errorf("application already exist: %s in %s", params.PackageName, params.TargetDomain)
	}

	applicationDir := filepath.Join("internal", "domain", params.TargetDomain, "application", params.PackageName)

	err = os.MkdirAll(applicationDir, 0o755)
	if err != nil {
		return "", fmt.Errorf("os: mkdir all: %w", err)
	}

	applicationFile, err := p.generateGoInitFile(
		ctx,
		applicationDir,
		params.StructName,
		params.PackageName,
		params.PortParam,
		TTApplication,
		params.AssertInterface,
	)
	if err != nil {
		err2 := os.RemoveAll(applicationDir)
		if err2 != nil {
			return "", fmt.Errorf("os: remove all: %w", err2)
		}
		return "", fmt.Errorf("generate application file: %w", err)
	}

	err = p.formatGoFiles(ctx, applicationFile)
	if err != nil {
		return "", fmt.Errorf("format go files: %w", err)
	}

	return applicationFile, nil
}

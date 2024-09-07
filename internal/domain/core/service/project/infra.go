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

func (p *Project) GetAllInfrastructes(ctx context.Context) ([]string, error) {
	infraPath := filepath.Join("internal", "infrastructure")

	infraCandidatePaths, err := filepath.Glob(filepath.Join(infraPath, "*"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	infraPaths := lo.Filter(infraCandidatePaths, func(d string, _ int) bool {
		stat, err2 := os.Stat(d)
		return err2 == nil && stat.IsDir()
	})

	infras := lo.Map(infraPaths, func(d string, _ int) string {
		return filepath.Base(d)
	})

	return infras, nil
}

func (p *Project) isInfraExist(ctx context.Context, targetInfra string) error {
	infras, err := p.GetAllInfrastructes(ctx)
	if err != nil {
		return fmt.Errorf("get all infrastructures: %w", err)
	}

	if !slices.Contains(infras, targetInfra) {
		return fmt.Errorf("target infrastructure not found: %s", targetInfra)
	}

	return nil
}

func (p *Project) CreateInfrastructure(ctx context.Context, params dto.CreateInfraParams) (string, error) {
	if params.PackageName == "" {
		params.PackageName = strings.ToLower(params.StructName)
	}

	err := p.ValidateInstanceName(params.StructName)
	if err != nil {
		return "", fmt.Errorf("validate service name: %w", err)
	}

	err = p.ValidatePkgName(params.PackageName)
	if err != nil {
		return "", fmt.Errorf("validate pkg name: %w", err)
	}

	err = p.isInfraExist(ctx, params.PackageName)
	if err == nil {
		return "", fmt.Errorf("is infra exist: %w", dto.ErrAlreadyExist)
	}

	infraDir := filepath.Join("internal", "infrastructure", params.PackageName)

	err = os.MkdirAll(infraDir, 0o755)
	if err != nil {
		return "", fmt.Errorf("os: mkdir all: %w", err)
	}

	infraFile, err := p.generateGoInitFile(
		ctx,
		infraDir,
		params.StructName,
		params.PackageName,
		params.PortParam,
		TTInfra,
		params.AssertInterface,
	)
	if err != nil {
		err2 := os.RemoveAll(infraDir)
		if err2 != nil {
			return "", fmt.Errorf("os: remove all: %w", err2)
		}
		return "", fmt.Errorf("generate infrastructure file: %w", err)
	}

	err = p.formatGoFiles(ctx, infraFile)
	if err != nil {
		return "", fmt.Errorf("format go files: %w", err)
	}

	return infraFile, nil
}

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

func (p *Project) GetAllPackages(ctx context.Context, showGlobal bool) ([]string, error) {
	pkgLocation := filepath.Join("internal", "pkg")

	if showGlobal {
		pkgLocation = "pkg"
	}

	pkgCandidatePaths, err := filepath.Glob(filepath.Join(pkgLocation, "*"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	pkgPaths := lo.Filter(pkgCandidatePaths, func(d string, _ int) bool {
		stat, err2 := os.Stat(d)
		return err2 == nil && stat.IsDir()
	})

	pkgs := lo.Map(pkgPaths, func(d string, _ int) string {
		return filepath.Base(d)
	})

	return pkgs, nil
}

func (p *Project) isPkgExist(ctx context.Context, targetPkg string, isGlobal bool) error {
	pkgs, err := p.GetAllPackages(ctx, isGlobal)
	if err != nil {
		return fmt.Errorf("get all packages: %w", err)
	}

	if !slices.Contains(pkgs, targetPkg) {
		return fmt.Errorf("target package not found: %s", targetPkg)
	}

	return nil
}

func (p *Project) CreatePackage(ctx context.Context, params dto.CreatePackageParams) (string, error) {
	if params.PackageName == "" {
		params.PackageName = strings.ToLower(params.StructName)
	}

	err := p.ValidateInstanceName(params.StructName)
	if err != nil {
		return "", fmt.Errorf("validate package name: %w", err)
	}

	err = p.ValidatePkgName(params.PackageName)
	if err != nil {
		return "", fmt.Errorf("validate pkg name: %w", err)
	}

	err = p.isPkgExist(ctx, params.PackageName, params.IsGlobal)
	if err == nil {
		return "", fmt.Errorf("package already exist: %s", params.StructName)
	}

	packageDir := filepath.Join("internal", "pkg", params.PackageName)

	if params.IsGlobal {
		packageDir = filepath.Join("pkg", params.PackageName)
	}

	err = os.MkdirAll(packageDir, 0o755)
	if err != nil {
		return "", fmt.Errorf("os: mkdir all: %w", err)
	}

	packageFile, err := p.generateGoInitFile(
		ctx,
		packageDir,
		params.StructName,
		params.PackageName,
		params.PortParam,
		TTPackage,
		params.AssertInterface,
	)
	if err != nil {
		err2 := os.RemoveAll(packageDir)
		if err2 != nil {
			return "", fmt.Errorf("os: remove all: %w", err2)
		}
		return "", fmt.Errorf("generate package file: %w", err)
	}

	err = p.formatGoFiles(ctx, packageFile)
	if err != nil {
		return "", fmt.Errorf("format go files: %w", err)
	}

	return packageFile, nil
}

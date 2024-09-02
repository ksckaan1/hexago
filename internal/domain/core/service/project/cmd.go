package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/samber/lo"
)

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

func (p *Project) isEntryPointExist(ctx context.Context, targetEntryPoint string) error {
	entryPoints, err := p.GetAllEntryPoints(ctx)
	if err != nil {
		return fmt.Errorf("get all entry points: %w", err)
	}

	if !slices.Contains(entryPoints, targetEntryPoint) {
		return fmt.Errorf("target entry point not found: %s", targetEntryPoint)
	}

	return nil
}

func (p *Project) CreateEntryPoint(ctx context.Context, params dto.CreateEntryPointParams) (string, error) {
	err := p.ValidateEntryPointName(params.PackageName)
	if err != nil {
		return "", fmt.Errorf("validate entry point name: %w", err)
	}

	err = p.isEntryPointExist(ctx, params.PackageName)
	if err == nil {
		return "", fmt.Errorf("entry point already exist: %s", params.PackageName)
	}

	entryPointPath := filepath.Join("cmd", params.PackageName)

	err = os.MkdirAll(entryPointPath, 0o755)
	if err != nil {
		return "", fmt.Errorf("os: mkdir all: %w", err)
	}

	cmdFile, err := assets.ReadFile("assets/templates/cmd.tmpl")
	if err != nil {
		return "", fmt.Errorf("assets: read file: %w", err)
	}

	cmdFilePath := filepath.Join(entryPointPath, "main.go")

	err = os.WriteFile(cmdFilePath, cmdFile, 0o644)
	if err != nil {
		return "", fmt.Errorf("os: write file: %w", err)
	}

	err = p.formatGoFiles(ctx, cmdFilePath)
	if err != nil {
		return "", fmt.Errorf("format go files: %w", err)
	}

	return cmdFilePath, nil
}

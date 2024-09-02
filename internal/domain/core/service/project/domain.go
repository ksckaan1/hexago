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

func (p *Project) isDomainExist(ctx context.Context, targetDomain string) error {
	domains, err := p.GetAllDomains(ctx)
	if err != nil {
		return fmt.Errorf("get all domains: %w", err)
	}

	if !slices.Contains(domains, targetDomain) {
		return fmt.Errorf("target domain not found: %w (%s)", dto.ErrDomainNotFound, targetDomain)
	}

	return nil
}

func (p *Project) CreateDomain(ctx context.Context, params dto.CreateDomainParams) error {
	err := p.isDomainExist(ctx, params.DomainName)
	if err == nil {
		return fmt.Errorf("domain already exist: %s", params.DomainName)
	}

	domainPath := filepath.Join("internal", "domain", params.DomainName)

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

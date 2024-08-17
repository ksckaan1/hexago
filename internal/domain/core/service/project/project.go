package projectservice

import (
	"context"
	"fmt"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/samber/do"
)

var _ port.ProjectService = (*ProjectService)(nil)

type ProjectService struct {
}

func New(i *do.Injector) (port.ProjectService, error) {
	return &ProjectService{}, nil
}

func (p *ProjectService) InitNewProject(ctx context.Context, params dto.InitNewProjectParams) error {
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

package project

import (
	"context"
	"fmt"

	"github.com/ksckaan1/hexago/config"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type Project struct {
	cfg *config.Config
}

func New(cfg *config.Config) (*Project, error) {
	return &Project{
		cfg: cfg,
	}, nil
}

func (p *Project) InitNewProject(ctx context.Context, params model.InitNewProjectParams) error {
	projectPath, err := p.createProjectDir(params.ProjectDirectory)
	if err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}

	err = p.createHexagoConfigs(projectPath)
	if err != nil {
		return fmt.Errorf("create hexago configs: %w", err)
	}

	err = p.addGitignore(projectPath)
	if err != nil {
		return fmt.Errorf("add gitignore: %w", err)
	}

	if params.CreateModule {
		err = p.initGoModule(ctx, params.ModuleName)
		if err != nil {
			return fmt.Errorf("init go module: %w", err)
		}
	}

	err = p.createProjectSubDirs()
	if err != nil {
		return fmt.Errorf("create project dirs: %w", err)
	}

	return nil
}

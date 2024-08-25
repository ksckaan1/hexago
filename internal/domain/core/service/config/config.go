package config

import (
	"fmt"
	"os"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
	"github.com/ksckaan1/hexago/internal/domain/core/port"
	"github.com/samber/do"
	"gopkg.in/yaml.v3"
)

var _ port.ConfigService = (*Config)(nil)

type Config struct {
	store model.Config
}

func New(i *do.Injector) (port.ConfigService, error) {
	return &Config{}, nil
}

func (c *Config) Load(cfgPath string) error {
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return fmt.Errorf("config file not found: %s", cfgPath)
	}
	defer cfgFile.Close()

	err = yaml.NewDecoder(cfgFile).Decode(&c.store)
	if err != nil {
		return fmt.Errorf("yaml: decode: %w", err)
	}

	return nil
}

func (c *Config) GetServiceTemplate() string {
	if c.store.Templates.Service == "" {
		return "std"
	}
	return c.store.Templates.Service
}

func (c *Config) GetApplicationTemplate() string {
	if c.store.Templates.Application == "" {
		return "std"
	}
	return c.store.Templates.Application
}

func (c *Config) GetInfrastructureTemplate() string {
	if c.store.Templates.Infrastructure == "" {
		return "std"
	}
	return c.store.Templates.Infrastructure
}

func (c *Config) GetPackageTemplate() string {
	if c.store.Templates.Package == "" {
		return "std"
	}
	return c.store.Templates.Package
}

func (c *Config) GetRunner(runner string) (*model.Runner, error) {
	if c.store.Runners == nil {
		return nil, dto.ErrRunnerNotImplemented
	}
	v, ok := c.store.Runners[runner]
	if !ok {
		return nil, dto.ErrRunnerNotImplemented
	}
	return v, nil
}

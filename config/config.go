package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/ksckaan1/hexago/internal/customerrors"
)

type Config struct {
	store    store
	location string
}

func New(location string) (*Config, error) {
	return &Config{
		location: location,
	}, nil
}

func (c *Config) Load() error {
	cfgFile, err := os.Open(c.location)
	if err != nil {
		return fmt.Errorf("config file not found: %s", c.location)
	}
	defer func() {
		err = cfgFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

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

func (c *Config) GetRunner(runner string) (*Runner, error) {
	if c.store.Runners == nil {
		return nil, customerrors.ErrRunnerNotImplemented
	}
	v, ok := c.store.Runners[runner]
	if !ok {
		return nil, customerrors.ErrRunnerNotImplemented
	}
	return v, nil
}

package project

import (
	"context"
	"fmt"
	"path/filepath"
)

func (p *Project) GetAllPorts(ctx context.Context, targetDomain string) ([]string, error) {
	err := p.isDomainExist(ctx, targetDomain)
	if err != nil {
		return nil, fmt.Errorf("is domain exist: %w", err)
	}

	portsPath := filepath.Join("internal", "domain", targetDomain, "port")

	portFilePaths, err := filepath.Glob(filepath.Join(portsPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("filepath: glob: %w", err)
	}

	allPorts := make([]string, 0)
	for i := range portFilePaths {
		ports, err2 := p.parseInterfaces(portFilePaths[i])
		if err2 != nil {
			return nil, fmt.Errorf("parse interfaces: %w", err2)
		}

		allPorts = append(allPorts, ports...)
	}

	return allPorts, nil
}

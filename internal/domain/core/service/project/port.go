package project

import (
	"context"
	"fmt"
	"path/filepath"
)

func (p *Project) GetAllPorts(_ context.Context) ([]string, error) {
	portsPath := filepath.Join("internal", "port")

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

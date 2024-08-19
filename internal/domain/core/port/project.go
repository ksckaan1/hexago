package port

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
)

type ProjectService interface {
	InitNewProject(ctx context.Context, params dto.InitNewProjectParams) error
	GetAllDomains(ctx context.Context) ([]string, error)
	CreateDomain(ctx context.Context, targetDomain string) error
	GetAllServices(ctx context.Context, targetDomain string) ([]string, error)
	CreateService(ctx context.Context, targetDomain, serviceName, pkgName, portName string) (string, error)
	GetAllPorts(ctx context.Context, targetDomain string) ([]string, error)
	GetAllApplications(ctx context.Context, targetDomain string) ([]string, error)
	CreateApplication(ctx context.Context, targetDomain, applicationName, pkgName, portName string) (string, error)
	GetAllEntryPoints(ctx context.Context) ([]string, error)
	CreateEntryPoint(ctx context.Context, entryPointName string) (string, error)
	GetAllInfrastructes(ctx context.Context) ([]string, error)
	CreateInfrastructure(ctx context.Context, infraName, pkgName, portName string) (string, error)
}

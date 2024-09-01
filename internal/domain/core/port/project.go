package port

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
)

type ProjectService interface {
	InitNewProject(ctx context.Context, params dto.InitNewProjectParams) error

	// Domain
	GetAllDomains(ctx context.Context) ([]string, error)
	CreateDomain(ctx context.Context, params dto.CreateDomainParams) error

	// Service
	GetAllServices(ctx context.Context, targetDomain string) ([]string, error)
	CreateService(ctx context.Context, params dto.CreateServiceParams) (string, error)

	// Port
	GetAllPorts(ctx context.Context, targetDomain string) ([]string, error)

	// Application
	GetAllApplications(ctx context.Context, targetDomain string) ([]string, error)
	CreateApplication(ctx context.Context, params dto.CreateApplicationParams) (string, error)

	// Entry Point (cmd)
	GetAllEntryPoints(ctx context.Context) ([]string, error)
	CreateEntryPoint(ctx context.Context, params dto.CreateEntryPointParams) (string, error)

	// Infrastructure
	GetAllInfrastructes(ctx context.Context) ([]string, error)
	CreateInfrastructure(ctx context.Context, params dto.CreateInfraParams) (string, error)

	// Package
	GetAllPackages(ctx context.Context, showGlobal bool) ([]string, error)
	CreatePackage(ctx context.Context, params dto.CreatePackageParams) (string, error)

	// Runner
	Run(ctx context.Context, command string, envVars []string, verbose bool) error

	// Doctor
	Doctor(ctx context.Context) (*dto.DoctorResult, error)
}

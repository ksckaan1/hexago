package treecmd

import "context"

type ProjectService interface {
	GetModuleName(modulePath ...string) (string, error)
	GetAllEntryPoints(ctx context.Context) ([]string, error)
	GetAllDomains(ctx context.Context) ([]string, error)
	GetAllServices(ctx context.Context, targetDomain string) ([]string, error)
	GetAllApplications(ctx context.Context, targetDomain string) ([]string, error)
	GetAllPackages(ctx context.Context, showGlobal bool) ([]string, error)
	GetAllInfrastructures(ctx context.Context) ([]string, error)
	GetAllPorts(ctx context.Context) ([]string, error)
}

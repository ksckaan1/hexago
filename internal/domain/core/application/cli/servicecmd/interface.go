package servicecmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	GetAllDomains(ctx context.Context) ([]string, error)
	ValidateInstanceName(instanceName string) error
	ValidatePkgName(pkgName string) error
	GetAllPorts(ctx context.Context) ([]string, error)
	CreateService(ctx context.Context, params model.CreateServiceParams) (string, error)
	GetAllServices(ctx context.Context, targetDomain string) ([]string, error)
}

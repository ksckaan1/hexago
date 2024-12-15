package infracmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	ValidateInstanceName(instanceName string) error
	ValidatePkgName(pkgName string) error
	GetAllPorts(ctx context.Context) ([]string, error)
	CreateInfrastructure(ctx context.Context, params model.CreateInfraParams) (string, error)
	GetAllInfrastructures(ctx context.Context) ([]string, error)
}

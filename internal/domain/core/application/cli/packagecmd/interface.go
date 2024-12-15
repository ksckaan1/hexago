package packagecmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	ValidateInstanceName(instanceName string) error
	ValidatePkgName(pkgName string) error
	GetAllPorts(ctx context.Context) ([]string, error)
	CreatePackage(ctx context.Context, params model.CreatePackageParams) (string, error)
	GetAllPackages(ctx context.Context, showGlobal bool) ([]string, error)
}

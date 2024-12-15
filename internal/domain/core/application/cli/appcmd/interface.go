package appcmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	GetAllDomains(ctx context.Context) ([]string, error)
	ValidateInstanceName(instanceName string) error
	ValidatePkgName(pkgName string) error
	GetAllPorts(ctx context.Context) ([]string, error)
	CreateApplication(ctx context.Context, params model.CreateApplicationParams) (string, error)
	GetAllApplications(ctx context.Context, targetDomain string) ([]string, error)
}

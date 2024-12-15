package domaincmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	GetAllDomains(ctx context.Context) ([]string, error)
	ValidatePkgName(pkgName string) error
	CreateDomain(ctx context.Context, params model.CreateDomainParams) error
}

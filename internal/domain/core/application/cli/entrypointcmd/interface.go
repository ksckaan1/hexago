package entrypointcmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	ValidateEntryPointName(entryPointName string) error
	CreateEntryPoint(ctx context.Context, params model.CreateEntryPointParams) (string, error)
	GetAllEntryPoints(ctx context.Context) ([]string, error)
}

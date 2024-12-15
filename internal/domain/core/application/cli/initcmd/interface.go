package initcmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	InitNewProject(ctx context.Context, params model.InitNewProjectParams) error
	GetModuleName(modulePath ...string) (string, error)
}

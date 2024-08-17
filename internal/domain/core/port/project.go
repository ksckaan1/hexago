package port

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
)

type ProjectService interface {
	InitNewProject(ctx context.Context, params dto.InitNewProjectParams) error
}

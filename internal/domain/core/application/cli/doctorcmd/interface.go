package doctorcmd

import (
	"context"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

type ProjectService interface {
	Doctor(ctx context.Context) (*model.DoctorResult, error)
}

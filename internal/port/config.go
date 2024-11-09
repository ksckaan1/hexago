package port

import "github.com/ksckaan1/hexago/internal/domain/core/model"

type ConfigService interface {
	Load(cfgPath string) error
	GetServiceTemplate() string
	GetApplicationTemplate() string
	GetInfrastructureTemplate() string
	GetPackageTemplate() string
	GetRunner(runner string) (*model.Runner, error)
}

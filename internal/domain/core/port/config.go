package port

type ConfigService interface {
	Load(cfgPath string) error
	GetServiceTemplate() string
	GetApplicationTemplate() string
	GetInfrastructureTemplate() string
	GetPackageTemplate() string
}

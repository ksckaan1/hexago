package dto

type InitNewProjectParams struct {
	ProjectDirectory string
	ModuleName       string
	CreateModule     bool
}

type CreateDomainParams struct {
	DomainName string
}

type CreateServiceParams struct {
	TargetDomain    string
	StructName      string
	PackageName     string
	PortParam       string
	AssertInterface bool
}

type CreateApplicationParams struct {
	TargetDomain    string
	StructName      string
	PackageName     string
	PortParam       string
	AssertInterface bool
}

type CreateEntryPointParams struct {
	PackageName string
}

type CreateInfraParams struct {
	StructName      string
	PackageName     string
	PortParam       string
	AssertInterface bool
}

type CreatePackageParams struct {
	StructName      string
	PackageName     string
	PortParam       string
	AssertInterface bool
	IsGlobal        bool
}

type DoctorResult struct {
	OSResult   string
	GoResult   string
	ImplResult string
}

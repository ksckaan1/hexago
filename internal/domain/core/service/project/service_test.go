package project

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/service/config"
	"github.com/stretchr/testify/require"
)

func TestCreateService(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx    func() context.Context
		params dto.CreateServiceParams
	}
	type want struct {
		err             require.ErrorAssertionFunc
		serviceFilePath string
	}

	tests := []struct {
		name string
		in
		args
		want
	}{
		{
			name: "domain not found",
			in: in{
				preRun: func(p *Project) error {
					return nil
				},
			},
			args: args{
				ctx: context.Background,
				params: dto.CreateServiceParams{
					TargetDomain:    "core",
					StructName:      "MyService",
					PackageName:     "myservice",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				serviceFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, dto.ErrDomainNotFound)
				},
			},
		},
		{
			name: "valid 1",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: dto.CreateServiceParams{
					TargetDomain:    "core",
					StructName:      "MyService",
					PackageName:     "myservice",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				serviceFilePath: filepath.Join("internal", "domain", "core", "service", "myservice", "myservice.go"),
				err:             require.NoError,
			},
		},
		{
			name: "valid with port",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: dto.CreateServiceParams{
					TargetDomain:    "core",
					StructName:      "MyService",
					PackageName:     "myservice",
					PortParam:       "io.Writer",
					AssertInterface: true,
				},
			},
			want: want{
				serviceFilePath: filepath.Join("internal", "domain", "core", "service", "myservice", "myservice.go"),
				err:             require.NoError,
			},
		},
		{
			name: "invalid port",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: dto.CreateServiceParams{
					TargetDomain:    "core",
					StructName:      "MyService",
					PackageName:     "myservice",
					PortParam:       "notexisting",
					AssertInterface: true,
				},
			},
			want: want{
				serviceFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, dto.ErrInvalidPortName{PortName: "notexisting"})
				},
			},
		},
		{
			name: "already existing",
			in: in{
				preRun: func(p *Project) error {
					err := p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
					if err != nil {
						return err
					}
					_, err = p.CreateService(context.Background(), dto.CreateServiceParams{
						TargetDomain: "core",
						StructName:   "MyService",
						PackageName:  "myservice",
					})
					return err
				},
			},
			args: args{
				ctx: context.Background,
				params: dto.CreateServiceParams{
					TargetDomain: "core",
					StructName:   "MyService",
					PackageName:  "myservice",
				},
			},
			want: want{
				serviceFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, dto.ErrAlreadyExist)
				},
			},
		},
		{
			name: "invalid instance name",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: dto.CreateServiceParams{
					TargetDomain:    "core",
					StructName:      "myService",
					PackageName:     "myservice",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				serviceFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, dto.ErrInvalidInstanceName)
				},
			},
		},
		{
			name: "invalid folder name",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: dto.CreateServiceParams{
					TargetDomain:    "core",
					StructName:      "MyService",
					PackageName:     "my-service",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				serviceFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, dto.ErrInvalidPkgName)
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			serviceFile, err := projectService.CreateService(tt.args.ctx(), tt.args.params)
			tt.want.err(t, err)
			require.Equal(t, tt.want.serviceFilePath, serviceFile)
		})
	}
}

func TestGetAllServices(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx          func() context.Context
		targetDomain string
	}
	type want struct {
		err      require.ErrorAssertionFunc
		services []string
	}

	tests := []struct {
		name string
		in
		args
		want
	}{
		{
			name: "domain not found",
			in: in{
				preRun: func(p *Project) error {
					return nil
				},
			},
			args: args{
				ctx:          context.Background,
				targetDomain: "core",
			},
			want: want{
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, dto.ErrDomainNotFound)
				},
				services: nil,
			},
		},
		{
			name: "valid",
			in: in{
				preRun: func(p *Project) error {
					err := p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
					if err != nil {
						return err
					}

					for i := range 3 {
						_, err = p.CreateService(context.Background(), dto.CreateServiceParams{
							TargetDomain: "core",
							StructName:   fmt.Sprintf("Example%d", i),
							PackageName:  fmt.Sprintf("example%d", i),
						})
						if err != nil {
							return err
						}
					}

					return nil
				},
			},
			args: args{
				ctx:          context.Background,
				targetDomain: "core",
			},
			want: want{
				err:      require.NoError,
				services: []string{"example0", "example1", "example2"},
			},
		},
		{
			name: "valid empty list",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx:          context.Background,
				targetDomain: "core",
			},
			want: want{
				err:      require.NoError,
				services: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			services, err := projectService.GetAllServices(tt.args.ctx(), tt.args.targetDomain)
			tt.want.err(t, err)
			require.Equal(t, tt.want.services, services)
		})
	}
}

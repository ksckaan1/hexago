package project

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ksckaan1/hexago/config"
	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

func TestCreateApplication(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx    func() context.Context
		params model.CreateApplicationParams
	}
	type want struct {
		err         require.ErrorAssertionFunc
		appFilePath string
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
				params: model.CreateApplicationParams{
					TargetDomain:    "core",
					StructName:      "MyApp",
					PackageName:     "myapp",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				appFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, customerrors.ErrDomainNotFound)
				},
			},
		},
		{
			name: "valid 1",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreateApplicationParams{
					TargetDomain:    "core",
					StructName:      "MyApp",
					PackageName:     "myapp",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				appFilePath: filepath.Join("internal", "domain", "core", "application", "myapp", "myapp.go"),
				err:         require.NoError,
			},
		},
		{
			name: "valid with port",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreateApplicationParams{
					TargetDomain:    "core",
					StructName:      "MyApp",
					PackageName:     "myapp",
					PortParam:       "io.Writer",
					AssertInterface: true,
				},
			},
			want: want{
				appFilePath: filepath.Join("internal", "domain", "core", "application", "myapp", "myapp.go"),
				err:         require.NoError,
			},
		},
		{
			name: "invalid port",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreateApplicationParams{
					TargetDomain:    "core",
					StructName:      "MyApp",
					PackageName:     "myapp",
					PortParam:       "notexisting",
					AssertInterface: true,
				},
			},
			want: want{
				appFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, customerrors.ErrInvalidPortName{PortName: "notexisting"})
				},
			},
		},
		{
			name: "already existing",
			in: in{
				preRun: func(p *Project) error {
					err := p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
					if err != nil {
						return err
					}
					_, err = p.CreateApplication(context.Background(), model.CreateApplicationParams{
						TargetDomain: "core",
						StructName:   "MyApp",
						PackageName:  "myapp",
					})
					return err
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreateApplicationParams{
					TargetDomain: "core",
					StructName:   "MyApp",
					PackageName:  "myapp",
				},
			},
			want: want{
				appFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, customerrors.ErrAlreadyExist)
				},
			},
		},
		{
			name: "invalid instance name",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreateApplicationParams{
					TargetDomain:    "core",
					StructName:      "myApp",
					PackageName:     "myapp",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				appFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, customerrors.ErrInvalidInstanceName)
				},
			},
		},
		{
			name: "invalid folder name",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreateApplicationParams{
					TargetDomain:    "core",
					StructName:      "MyApp",
					PackageName:     "my-app",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				appFilePath: "",
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, customerrors.ErrInvalidPkgName)
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

			appFile, err := projectService.CreateApplication(tt.args.ctx(), tt.args.params)
			tt.want.err(t, err)
			require.Equal(t, tt.want.appFilePath, appFile)
		})
	}
}

func TestGetAllApplications(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx          func() context.Context
		targetDomain string
	}
	type want struct {
		err  require.ErrorAssertionFunc
		apps []string
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
					require.ErrorIs(tt, err, customerrors.ErrDomainNotFound)
				},
				apps: nil,
			},
		},
		{
			name: "valid",
			in: in{
				preRun: func(p *Project) error {
					err := p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
					if err != nil {
						return err
					}

					for i := range 3 {
						_, err = p.CreateApplication(context.Background(), model.CreateApplicationParams{
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
				err:  require.NoError,
				apps: []string{"example0", "example1", "example2"},
			},
		},
		{
			name: "valid empty list",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(context.Background(), model.InitNewProjectParams{
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
				err:  require.NoError,
				apps: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			apps, err := projectService.GetAllApplications(tt.args.ctx(), tt.args.targetDomain)
			tt.want.err(t, err)
			require.Equal(t, tt.want.apps, apps)
		})
	}
}

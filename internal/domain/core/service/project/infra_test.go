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

func TestCreateInfrastructure(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx    func() context.Context
		params model.CreateInfraParams
	}
	type want struct {
		err           require.ErrorAssertionFunc
		infraFilePath string
	}

	tests := []struct {
		name string
		in
		args
		want
	}{
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
				params: model.CreateInfraParams{
					StructName:      "MyInfra",
					PackageName:     "myinfra",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				infraFilePath: filepath.Join("internal", "infrastructure", "myinfra", "myinfra.go"),
				err:           require.NoError,
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
				params: model.CreateInfraParams{
					StructName:      "MyInfra",
					PackageName:     "myinfra",
					PortParam:       "io.Writer",
					AssertInterface: true,
				},
			},
			want: want{
				infraFilePath: filepath.Join("internal", "infrastructure", "myinfra", "myinfra.go"),
				err:           require.NoError,
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
				params: model.CreateInfraParams{
					StructName:      "MyApp",
					PackageName:     "myapp",
					PortParam:       "notexisting",
					AssertInterface: true,
				},
			},
			want: want{
				infraFilePath: "",
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
					_, err = p.CreateInfrastructure(context.Background(), model.CreateInfraParams{
						StructName:  "MyApp",
						PackageName: "myapp",
					})
					return err
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreateInfraParams{
					StructName:  "MyApp",
					PackageName: "myapp",
				},
			},
			want: want{
				infraFilePath: "",
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
				params: model.CreateInfraParams{
					StructName:      "myApp",
					PackageName:     "myapp",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				infraFilePath: "",
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
				params: model.CreateInfraParams{
					StructName:      "MyApp",
					PackageName:     "my-app",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				infraFilePath: "",
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

			infraFile, err := projectService.CreateInfrastructure(tt.args.ctx(), tt.args.params)
			tt.want.err(t, err)
			require.Equal(t, tt.want.infraFilePath, infraFile)
		})
	}
}

func TestGetAllInfrastructures(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx func() context.Context
	}
	type want struct {
		err             require.ErrorAssertionFunc
		infrastructures []string
	}

	tests := []struct {
		name string
		in
		args
		want
	}{
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
						_, err = p.CreateInfrastructure(context.Background(), model.CreateInfraParams{
							StructName:  fmt.Sprintf("Example%d", i),
							PackageName: fmt.Sprintf("example%d", i),
						})
						if err != nil {
							return err
						}
					}

					return nil
				},
			},
			args: args{
				ctx: context.Background,
			},
			want: want{
				err:             require.NoError,
				infrastructures: []string{"example0", "example1", "example2"},
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
				ctx: context.Background,
			},
			want: want{
				err:             require.NoError,
				infrastructures: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			infrastructures, err := projectService.GetAllInfrastructures(tt.args.ctx())
			tt.want.err(t, err)
			require.Equal(t, tt.want.infrastructures, infrastructures)
		})
	}
}

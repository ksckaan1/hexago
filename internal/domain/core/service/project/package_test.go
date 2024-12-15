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

func TestCreatePackage(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx    func() context.Context
		params model.CreatePackageParams
	}
	type want struct {
		err             require.ErrorAssertionFunc
		packageFilePath string
	}

	tests := []struct {
		name string
		in
		args
		want
	}{
		{
			name: "valid internal",
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
				params: model.CreatePackageParams{
					StructName:      "MyPkg",
					PackageName:     "mypkg",
					IsGlobal:        false,
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				packageFilePath: filepath.Join("internal", "pkg", "mypkg", "mypkg.go"),
				err:             require.NoError,
			},
		},
		{
			name: "valid global",
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
				params: model.CreatePackageParams{
					StructName:      "MyPkg",
					PackageName:     "mypkg",
					IsGlobal:        true,
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				packageFilePath: filepath.Join("pkg", "mypkg", "mypkg.go"),
				err:             require.NoError,
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
				params: model.CreatePackageParams{
					StructName:      "MyPkg",
					PackageName:     "mypkg",
					PortParam:       "io.Writer",
					AssertInterface: true,
				},
			},
			want: want{
				packageFilePath: filepath.Join("internal", "pkg", "mypkg", "mypkg.go"),
				err:             require.NoError,
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
				params: model.CreatePackageParams{
					StructName:      "MyPkg",
					PackageName:     "mypkg",
					PortParam:       "notexisting",
					AssertInterface: true,
				},
			},
			want: want{
				packageFilePath: "",
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
					_, err = p.CreatePackage(context.Background(), model.CreatePackageParams{
						StructName:  "MyPkg",
						PackageName: "mypkg",
					})
					return err
				},
			},
			args: args{
				ctx: context.Background,
				params: model.CreatePackageParams{
					StructName:  "MyPkg",
					PackageName: "mypkg",
				},
			},
			want: want{
				packageFilePath: "",
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
				params: model.CreatePackageParams{
					StructName:      "myPkg",
					PackageName:     "mypkg",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				packageFilePath: "",
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
				params: model.CreatePackageParams{
					StructName:      "MyPkg",
					PackageName:     "my-pkg",
					PortParam:       "",
					AssertInterface: false,
				},
			},
			want: want{
				packageFilePath: "",
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

			packageFile, err := projectService.CreatePackage(tt.args.ctx(), tt.args.params)
			tt.want.err(t, err)
			require.Equal(t, tt.want.packageFilePath, packageFile)
		})
	}
}

func TestGetAllPackages(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx func() context.Context
	}
	type want struct {
		err      require.ErrorAssertionFunc
		packages []string
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
						_, err = p.CreatePackage(context.Background(), model.CreatePackageParams{
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
				err:      require.NoError,
				packages: []string{"example0", "example1", "example2"},
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
				err:      require.NoError,
				packages: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			packages, err := projectService.GetAllPackages(tt.args.ctx(), false)
			tt.want.err(t, err)
			require.Equal(t, tt.want.packages, packages)
		})
	}
}

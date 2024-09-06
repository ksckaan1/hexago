package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/ksckaan1/hexago/internal/domain/core/service/config"
	"github.com/stretchr/testify/require"
)

func TestCreateDomain(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx        func() context.Context
		domainName string
	}
	type want struct {
		err require.ErrorAssertionFunc
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
					return p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: t.TempDir(),
						ModuleName:       "my-project",
						CreateModule:     true,
					})
				},
			},
			args: args{
				ctx:        context.Background,
				domainName: "example",
			},
			want: want{
				err: require.NoError,
			},
		},
		{
			name: "already exist",
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
				ctx:        context.Background,
				domainName: "core",
			},
			want: want{
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(tt, err, dto.ErrAlreadyExist)
				},
			},
		},
		{
			name: "invalid domain name",
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
				ctx:        context.Background,
				domainName: "invalid domain name",
			},
			want: want{
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

			err := projectService.CreateDomain(tt.args.ctx(), dto.CreateDomainParams{
				DomainName: tt.args.domainName,
			})
			tt.want.err(t, err)
		})
	}
}

func TestGetAllDomains(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx func() context.Context
	}
	type want struct {
		domains []string
		err     require.ErrorAssertionFunc
	}

	tests := []struct {
		name string
		in
		args
		want
	}{
		{
			name: "init state",
			in: in{
				preRun: func(p *Project) error {
					return p.InitNewProject(
						context.Background(),
						dto.InitNewProjectParams{
							ProjectDirectory: t.TempDir(),
							ModuleName:       "my-project",
						},
					)
				},
			},
			args: args{
				ctx: context.Background,
			},
			want: want{
				domains: []string{"core"},
				err:     require.NoError,
			},
		},
		{
			name: "empty list",
			in: in{
				preRun: func(p *Project) error {
					dir := t.TempDir()
					err := p.InitNewProject(
						context.Background(),
						dto.InitNewProjectParams{
							ProjectDirectory: dir,
							ModuleName:       "my-project",
						},
					)
					if err != nil {
						return err
					}
					return os.RemoveAll(filepath.Join(dir, "internal", "domain", "core"))
				},
			},
			args: args{
				ctx: context.Background,
			},
			want: want{
				domains: []string{},
				err:     require.NoError,
			},
		},
		{
			name: "multi domain",
			in: in{
				preRun: func(p *Project) error {
					dir := t.TempDir()
					err := p.InitNewProject(
						context.Background(),
						dto.InitNewProjectParams{
							ProjectDirectory: dir,
							ModuleName:       "my-project",
						},
					)
					if err != nil {
						return err
					}

					for i := range 2 {
						err = p.CreateDomain(context.Background(), dto.CreateDomainParams{
							DomainName: fmt.Sprintf("domain%d", i),
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
				domains: []string{"core", "domain0", "domain1"},
				err:     require.NoError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			domains, err := projectService.GetAllDomains(tt.args.ctx())
			tt.want.err(t, err)
			require.Equal(t, tt.want.domains, domains)
		})
	}
}

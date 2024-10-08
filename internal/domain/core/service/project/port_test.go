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

func TestGetAllPorts(t *testing.T) {
	type in struct {
		preRun func(p *Project) error
	}
	type args struct {
		ctx          func() context.Context
		targetDomain string
	}
	type want struct {
		err   require.ErrorAssertionFunc
		ports []string
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
					dir := t.TempDir()
					err := p.InitNewProject(context.Background(), dto.InitNewProjectParams{
						ProjectDirectory: dir,
						ModuleName:       "my-project",
						CreateModule:     true,
					})
					if err != nil {
						return err
					}

					for i := range 3 {
						err = os.WriteFile(filepath.Join(dir, "internal", "domain", "core", "port", fmt.Sprintf("example%d.go", i)), []byte(fmt.Sprintf("package port\ntype Example%d interface {}\n", i)), 0o644)
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
				err:   require.NoError,
				ports: []string{"Example0", "Example1", "Example2"},
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
				err:   require.NoError,
				ports: []string{},
			},
		},
		{
			name: "domain not found",
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
				targetDomain: "notexisting",
			},
			want: want{
				err: func(tt require.TestingT, err error, i ...interface{}) {
					require.ErrorIs(t, err, dto.ErrDomainNotFound)
				},
				ports: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			ports, err := projectService.GetAllPorts(tt.args.ctx(), tt.args.targetDomain)
			tt.want.err(t, err)
			require.Equal(t, tt.want.ports, ports)
		})
	}
}

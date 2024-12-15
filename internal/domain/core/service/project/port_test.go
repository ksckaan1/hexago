package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ksckaan1/hexago/config"
	"github.com/ksckaan1/hexago/internal/domain/core/model"
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
					err := p.InitNewProject(context.Background(), model.InitNewProjectParams{
						ProjectDirectory: dir,
						ModuleName:       "my-project",
						CreateModule:     true,
					})
					if err != nil {
						return err
					}

					for i := range 3 {
						err = os.WriteFile(filepath.Join(dir, "internal", "port", fmt.Sprintf("example%d.go", i)), []byte(fmt.Sprintf("package port\ntype Example%d interface {}\n", i)), 0o644)
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
				err:   require.NoError,
				ports: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{
				cfg: &config.Config{},
			}
			require.NoError(t, tt.in.preRun(projectService))

			ports, err := projectService.GetAllPorts(tt.args.ctx())
			tt.want.err(t, err)
			require.Equal(t, tt.want.ports, ports)
		})
	}
}

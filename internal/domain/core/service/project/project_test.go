package projectservice

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/stretchr/testify/require"
)

func TestInitNewProject(t *testing.T) {
	type args struct {
		ctx func() context.Context
	}
	type want struct {
		err require.ErrorAssertionFunc
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "valid1",
			args: args{
				ctx: context.Background,
			},
			want: want{
				err: require.NoError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &ProjectService{}

			projectDir := filepath.Join(t.TempDir(), "my-project")

			err := projectService.InitNewProject(
				tt.args.ctx(),
				dto.InitNewProjectParams{
					ProjectDirectory: projectDir,
					ModuleName:       "my-project",
				},
			)

			tt.want.err(t, err)
		})
	}
}

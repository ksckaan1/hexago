package project

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/stretchr/testify/require"
)

func TestInitNewProject(t *testing.T) {
	type args struct {
		ctx           func() context.Context
		projectFolder string
		moduleName    string
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
			name: "valid dir",
			args: args{
				ctx:           context.Background,
				projectFolder: "my-project",
				moduleName:    "my-module",
			},
			want: want{
				err: require.NoError,
			},
		},
		{
			name: "dot dir",
			args: args{
				ctx:           context.Background,
				projectFolder: ".",
				moduleName:    "my-module",
			},
			want: want{
				err: require.NoError,
			},
		},
		{
			name: "empty module name",
			args: args{
				ctx:           context.Background,
				projectFolder: ".",
				moduleName:    "",
			},
			want: want{
				err: require.NoError,
			},
		},
		{
			name: "empty folder name",
			args: args{
				ctx:           context.Background,
				projectFolder: "",
				moduleName:    "",
			},
			want: want{
				err: require.NoError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{}

			projectDir := filepath.Join(t.TempDir(), tt.args.projectFolder)

			err := projectService.InitNewProject(
				tt.args.ctx(),
				dto.InitNewProjectParams{
					ProjectDirectory: projectDir,
					ModuleName:       tt.args.moduleName,
				},
			)

			tt.want.err(t, err)
		})
	}
}

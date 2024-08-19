package project

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
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
			projectService := &Project{}

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

func TestGetAllDomains(t *testing.T) {
	type args struct {
		ctx func() context.Context
	}
	type want struct {
		domains []string
		err     require.ErrorAssertionFunc
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
				domains: []string{"core"},
				err:     require.NoError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectService := &Project{}

			projectDir := filepath.Join(t.TempDir(), "my-project")

			err := projectService.InitNewProject(
				tt.args.ctx(),
				dto.InitNewProjectParams{
					ProjectDirectory: projectDir,
					ModuleName:       "my-project",
				},
			)
			require.NoError(t, err)

			domains, err := projectService.GetAllDomains(tt.args.ctx())
			tt.want.err(t, err)
			require.Equal(t, tt.want.domains, domains)
		})
	}
}

func TestAssets(t *testing.T) {
	content, err := assets.ReadFile("assets/config.yaml")
	require.NoError(t, err)

	t.Log(string(content))
}

func TestParseInterfaces(t *testing.T) {
	data := `
type UserService interface{
	Example() error
}

// type UserRepository interface{
	Example() error
}

type Company interface {
	Example() error
}`

	submatches := rgxInterfaces.FindAllStringSubmatch(data, -1)
	for i := range submatches {
		for j := range submatches[i] {
			t.Log(i, j, submatches[i][j])
		}
	}
}

func TestParseGoMod(t *testing.T) {
	f, err := os.Open("../../../../../go.mod")
	require.NoError(t, err)
	defer f.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, f)
	require.NoError(t, err)

	modFile, err := modfile.Parse("go.mod", buf.Bytes(), nil)
	require.NoError(t, err)

	t.Log(modFile.Module.Mod.Path)
}

func TestParseDefaultModuleName(t *testing.T) {
	defaultName := "."

	if defaultName == "." {
		abs, err := filepath.Abs(defaultName)
		if err != nil {
			require.NoError(t, err)
		}

		defaultName = filepath.Base(abs)
	}

	t.Log(defaultName)
}

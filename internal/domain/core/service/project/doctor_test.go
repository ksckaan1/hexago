package project

import (
	"context"
	"testing"

	"github.com/ksckaan1/hexago/internal/domain/core/service/config"
	"github.com/stretchr/testify/require"
)

func TestDoctor(t *testing.T) {
	projectService := &Project{
		cfg: &config.Config{},
	}

	result, err := projectService.Doctor(context.Background())
	require.NoError(t, err)

	require.NotEmpty(t, result.OSResult)
	require.NotEmpty(t, result.GoResult)
	require.NotEmpty(t, result.ImplResult)
}

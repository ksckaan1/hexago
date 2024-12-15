//go:build integration
// +build integration

package project

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ksckaan1/hexago/config"
)

func TestDoctor(t *testing.T) {
	projectService := &Project{
		cfg: &config.Config{},
	}

	result, err := projectService.Doctor(context.Background())
	require.NoError(t, err)

	require.NotEmpty(t, result.OSResult)
	require.True(t, result.GoResult.IsInstalled, result.GoResult.Output)
	require.True(t, result.ImplResult.IsInstalled, result.ImplResult.Output)
}

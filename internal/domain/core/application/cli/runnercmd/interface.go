package runnercmd

import "context"

type ProjectService interface {
	Run(ctx context.Context, command string, envVars []string, verbose bool) error
}

package portcmd

import "context"

type ProjectService interface {
	GetAllPorts(ctx context.Context) ([]string, error)
}

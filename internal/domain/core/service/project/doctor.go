package project

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
)

func (p *Project) Doctor(ctx context.Context) (*dto.DoctorResult, error) {
	goCommand, err := p.isToolInstalled(ctx, "go version")
	if err != nil {
		return nil, fmt.Errorf("check command: %w", err)
	}

	implCommand, err := p.isToolInstalled(ctx, "impl Murmur hash.Hash")
	if err != nil {
		return nil, fmt.Errorf("check command: %w", err)
	}

	return &dto.DoctorResult{
		OSResult:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GoResult:   goCommand,
		ImplResult: implCommand,
	}, nil
}

func (p *Project) isToolInstalled(ctx context.Context, command string) (dto.Tool, error) {
	term, err := p.getTerminalName()
	if err != nil {
		return dto.Tool{}, fmt.Errorf("get terminal name: %w", err)
	}

	cmd := exec.CommandContext(ctx, term, "-c", command)
	cmd.Env = os.Environ()
	stdOut, stdErr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdOut, stdErr

	err = cmd.Run()
	if err != nil {
		return dto.Tool{
			IsInstalled: false,
			Output:      stdErr.String(),
		}, nil
	}

	return dto.Tool{
		IsInstalled: true,
		Output:      stdOut.String(),
	}, nil
}

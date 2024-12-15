package project

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

func (p *Project) Doctor(ctx context.Context) (*model.DoctorResult, error) {
	goCommand, err := p.isToolInstalled(ctx, "go version")
	if err != nil {
		return nil, fmt.Errorf("check command: %w", err)
	}

	implCommand, err := p.isToolInstalled(ctx, "impl Murmur hash.Hash")
	if err != nil {
		return nil, fmt.Errorf("check command: %w", err)
	}

	return &model.DoctorResult{
		OSResult:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GoResult:   goCommand,
		ImplResult: implCommand,
	}, nil
}

func (p *Project) isToolInstalled(ctx context.Context, command string) (model.Tool, error) {
	term, err := p.getTerminalName()
	if err != nil {
		return model.Tool{}, fmt.Errorf("get terminal name: %w", err)
	}

	cmd := exec.CommandContext(ctx, term, "-c", command)
	cmd.Env = os.Environ()
	stdOut, stdErr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdOut, stdErr

	err = cmd.Run()
	if err != nil {
		return model.Tool{
			IsInstalled: false,
			Output:      stdErr.String(),
		}, nil
	}

	return model.Tool{
		IsInstalled: true,
		Output:      strings.TrimSpace(stdOut.String()),
	}, nil
}

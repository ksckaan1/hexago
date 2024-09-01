package project

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
)

func (p *Project) Doctor(ctx context.Context) (*dto.DoctorResult, error) {
	goCommand, err := p.checkCommand(ctx, "go version")
	if err != nil {
		return nil, fmt.Errorf("check command: %w", err)
	}

	implCommand, err := p.checkCommand(ctx, "impl Murmur hash.Hash")
	if err != nil {
		return nil, fmt.Errorf("check command: %w", err)
	}

	if implCommand != "" {
		implCommand = "installed"
	}

	return &dto.DoctorResult{
		OSResult:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GoResult:   goCommand,
		ImplResult: implCommand,
	}, nil
}

func (p *Project) checkCommand(ctx context.Context, command string) (string, error) {
	term, err := p.getTerminalName(ctx)
	if err != nil {
		return "", fmt.Errorf("get terminal name: %w", err)
	}

	cmd := exec.CommandContext(ctx, term, "-c", command)
	cmd.Env = os.Environ()
	stdOut := &bytes.Buffer{}
	cmd.Stdout = stdOut

	err = cmd.Run()
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(stdOut.String()), nil
}

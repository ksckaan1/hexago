package project

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"golang.org/x/mod/modfile"
)

func (*Project) initGoModule(ctx context.Context, moduleName string) error {
	_, err := os.Stat("go.mod")
	if !os.IsNotExist(err) {
		err = os.Remove("go.mod")
		if err != nil {
			return fmt.Errorf("os: remove: %w", err)
		}
	}

	if moduleName == "." {
		abs, err := filepath.Abs(moduleName)
		if err != nil {
			return fmt.Errorf("filepath: abs: %w", err)
		}

		moduleName = filepath.Base(abs)
	}

	cmd := exec.CommandContext(ctx, "go", "mod", "init", moduleName)
	stdErr := &bytes.Buffer{}
	cmd.Stderr = stdErr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd: run: %w", dto.ErrInitGoModule{Message: strings.TrimSpace(stdErr.String())})
	}

	return nil
}

func (p *Project) GetModuleName(modulePath ...string) (string, error) {
	mp := "go.mod"
	if len(modulePath) > 0 {
		mp = modulePath[0]
	}

	f, err := os.Open(mp)
	if err != nil {
		return "", fmt.Errorf("module file not found: %s", mp)
	}
	defer f.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, f)
	if err != nil {
		return "", fmt.Errorf("io: copy: %w", err)
	}

	modFile, err := modfile.Parse("go.mod", buf.Bytes(), nil)
	if err != nil {
		return "", fmt.Errorf("modfile parse: %w", err)
	}

	return modFile.Module.Mod.Path, nil
}

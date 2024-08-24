package project

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

func (*Project) initGoModule(ctx context.Context, moduleName string) error {
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

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd: run: %w", errors.New(strings.TrimSpace(stdErr.String())))
	}

	return nil
}

func (p *Project) getModuleName() (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("module file not found: go.mod")
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

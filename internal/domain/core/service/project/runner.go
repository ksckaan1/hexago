package project

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/ksckaan1/hexago/internal/domain/core/model"
)

func (p *Project) Run(ctx context.Context, commandName string, envVars []string, verbose bool) error {
	runner, err := p.getCommandInfo(ctx, commandName)
	if err != nil {
		return fmt.Errorf("get command info: %w", err)
	}

	envs := slices.Concat(os.Environ(), envVars, runner.EnvVars)

	if verbose {
		fmt.Println("cmd:", runner.Cmd)
		fmt.Println("env:", slices.Concat(envVars, runner.EnvVars))
		fmt.Println("log.disabled:", runner.Log.Disabled)
		fmt.Println("log.seperate_files:", runner.Log.SeperateFiles)
		fmt.Println("log.overwrite:", runner.Log.Overwrite)
		fmt.Println("-----------------------------------")
	}

	err = p.runCmd(ctx, commandName, runner, envs)
	if err != nil {
		return fmt.Errorf("run cmd: %w", err)
	}

	return nil
}

func (p *Project) getCommandInfo(ctx context.Context, commandName string) (*model.Runner, error) {
	runner, err := p.cfg.GetRunner(commandName)
	if err == nil && runner.Cmd != "" {
		return runner, nil
	}

	if runner == nil {
		runner = &model.Runner{}
	}

	entryPoints, err := p.GetAllEntryPoints(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all entry points: %w", err)
	}

	if !slices.Contains(entryPoints, commandName) {
		return nil, fmt.Errorf("entry point or runner not found")
	}

	runner.Cmd = "go run ./cmd/" + commandName

	return runner, nil
}

func (p *Project) prepareLogFiles(commandName string, runner *model.Runner, cmd *exec.Cmd) ([]io.Closer, error) {
	if runner.Log.Disabled {
		cmd.Stderr, cmd.Stdout = os.Stderr, os.Stdout
		return nil, nil
	}

	err := os.MkdirAll("logs", 0o755)
	if err != nil {
		return nil, fmt.Errorf("os: mkdir all: %w", err)
	}

	closers := make([]io.Closer, 0)

	if runner.Log.SeperateFiles {
		stdErrFile, err2 := p.createLogFile(commandName+".stderr", runner.Log.Overwrite)
		if err2 != nil {
			return nil, fmt.Errorf("create log file: %w", err2)
		}
		closers = append(closers, stdErrFile)

		cmd.Stderr = io.MultiWriter(os.Stderr, stdErrFile)

		stdOutFile, err2 := p.createLogFile(commandName+".stdout", runner.Log.Overwrite)
		if err2 != nil {
			return nil, fmt.Errorf("create log file: %w", err2)
		}
		closers = append(closers, stdOutFile)

		cmd.Stdout = io.MultiWriter(os.Stdout, stdOutFile)
	} else {
		logFile, err2 := p.createLogFile(commandName, runner.Log.Overwrite)
		if err2 != nil {
			return nil, fmt.Errorf("create log file: %w", err2)
		}
		closers = append(closers, logFile)

		cmd.Stderr = io.MultiWriter(os.Stderr, logFile)
		cmd.Stdout = io.MultiWriter(os.Stdout, logFile)
	}

	return closers, nil
}

func (p *Project) createLogFile(name string, overwrite bool) (*os.File, error) {
	filePath := filepath.Join("logs", fmt.Sprintf("%s.log", name))
	if !overwrite {
		logFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			return logFile, nil
		}
	}
	logFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("os: create: %w", err)
	}
	return logFile, nil
}

func (p *Project) runCmd(ctx context.Context, commandName string, runner *model.Runner, envs []string) error {
	term, err := p.getTerminalName()
	if err != nil {
		return fmt.Errorf("get terminal name: %w", err)
	}

	cmd := exec.CommandContext(ctx, term, "-c", runner.Cmd)
	cmd.Env = envs
	cmd.Stdin = os.Stdin

	closers, err := p.prepareLogFiles(commandName, runner, cmd)
	if err != nil {
		return fmt.Errorf("prepare log files: %w", err)
	}
	defer func() {
		for i := range closers {
			closers[i].Close()
		}
	}()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd: run: %w", err)
	}

	return nil
}

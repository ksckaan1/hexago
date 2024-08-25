package project

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"slices"
)

func (p *Project) Run(ctx context.Context, commandName string, envVars []string, verbose bool) error {
	commandInfo, err := p.getCommandInfo(ctx, commandName)
	if err != nil {
		return fmt.Errorf("get command info: %w", err)
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", commandInfo.Command)
	cmd.Env = slices.Concat(os.Environ(), envVars, commandInfo.EnvVars)
	cmd.Stdin = os.Stdin

	if !commandInfo.DisableLogFiles {
		err = os.MkdirAll("logs", 0o755)
		if err != nil {
			return fmt.Errorf("os: mkdir all: %w", err)
		}

		if commandInfo.SeperateLogFiles {
			stdErrFile, err := os.Create(fmt.Sprintf("logs/%s.stderr.log", commandName))
			if err != nil {
				log.Fatalln(err)
			}
			defer stdErrFile.Close()

			cmd.Stderr = io.MultiWriter(os.Stderr, stdErrFile)

			stdOutFile, err := os.Create(fmt.Sprintf("logs/%s.stdout.log", commandName))
			if err != nil {
				log.Fatalln(err)
			}
			defer stdOutFile.Close()

			cmd.Stdout = io.MultiWriter(os.Stdout, stdOutFile)
		} else {
			logFile, err := os.Create(fmt.Sprintf("logs/%s.log", commandName))
			if err != nil {
				log.Fatalln(err)
			}
			defer logFile.Close()

			cmd.Stderr = io.MultiWriter(os.Stderr, logFile)
			cmd.Stdout = io.MultiWriter(os.Stdout, logFile)
		}
	} else {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	sigCh := make(chan os.Signal, 1)
	go func() {
		sig := <-sigCh
		fmt.Println("signal received:", sig.String())
		cmd.Process.Signal(sig)
	}()
	signal.Notify(sigCh, os.Kill, os.Interrupt)

	if verbose {
		fmt.Println("command:", cmd.String())
		fmt.Println("env vars:", slices.Concat(envVars, commandInfo.EnvVars))
		fmt.Println("seperate log files:", commandInfo.SeperateLogFiles)
		fmt.Println("disable log files:", commandInfo.DisableLogFiles)
	}

	err = cmd.Run()
	if err != nil {
		os.Exit(cmd.ProcessState.ExitCode())
	}

	return nil
}

type CommandInfo struct {
	Command          string
	EnvVars          []string
	SeperateLogFiles bool
	DisableLogFiles  bool
}

func (p *Project) getCommandInfo(ctx context.Context, commandName string) (*CommandInfo, error) {
	cmdInfo := &CommandInfo{}

	runner, err := p.cfg.GetRunner(commandName)
	if err == nil {
		cmdInfo.Command = runner.Cmd
		cmdInfo.EnvVars = runner.EnvVars
		cmdInfo.SeperateLogFiles = runner.SeperateLogFiles
		cmdInfo.DisableLogFiles = runner.DisableLogFiles
	}

	if cmdInfo.Command != "" {
		return cmdInfo, nil
	}

	entryPoints, err := p.GetAllEntryPoints(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all entry points: %w", err)
	}

	if !slices.Contains(entryPoints, commandName) {
		return nil, fmt.Errorf("entry point or runner not found")
	}

	cmdInfo.Command = "go run ./cmd/" + commandName

	return cmdInfo, nil
}

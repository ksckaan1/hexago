package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ksckaan1/hexago/config"
	"github.com/ksckaan1/hexago/internal/customerrors"
	"github.com/ksckaan1/hexago/internal/domain/core/service/project"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
)

func main() {
	ctx := context.Background()

	cfgLocation := cmp.Or(
		os.Getenv("HEXAGO_CONFIG"),
		".hexago/config.yaml",
	)

	cfg, err := config.New(cfgLocation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tl, err := tuilog.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	projectService, err := project.New(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app, err := initCommands(projectService, tl, cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = app.Run(ctx)
	if err != nil {
		if !errors.Is(err, customerrors.ErrSuppressed) {
			fmt.Println(unwrapAllErrors(err))
		}
		os.Exit(1)
	}
}

func unwrapAllErrors(err error) error {
	if err == nil {
		return nil
	}
	for {
		uw := errors.Unwrap(err)
		if uw == nil {
			return err
		}
		err = uw
	}
}

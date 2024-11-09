package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ksckaan1/hexago/internal/domain/core/dto"
	"os"

	"github.com/ksckaan1/hexago/internal/domain/core/application/cli"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/appcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/doctorcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/domaincmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/entrypointcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/infracmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/packagecmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/portcmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/runnercmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/servicecmd"
	"github.com/ksckaan1/hexago/internal/domain/core/application/cli/treecmd"
	"github.com/ksckaan1/hexago/internal/domain/core/service/config"
	projectservice "github.com/ksckaan1/hexago/internal/domain/core/service/project"
	"github.com/ksckaan1/hexago/internal/pkg/tuilog"
	"github.com/samber/do"
)

func main() {
	ctx := context.Background()

	i := do.New()
	do.Provide(i, tuilog.New)
	do.Provide(i, cli.New)
	do.Provide(i, cli.NewRootCommand)
	do.Provide(i, cli.NewInitCommand)
	do.Provide(i, domaincmd.NewDomainCommand)
	do.Provide(i, domaincmd.NewDomainLSCommand)
	do.Provide(i, domaincmd.NewDomainCreateCommand)
	do.Provide(i, servicecmd.NewServiceCommand)
	do.Provide(i, servicecmd.NewServiceLSCommand)
	do.Provide(i, servicecmd.NewServiceCreateCommand)
	do.Provide(i, portcmd.NewPortCommand)
	do.Provide(i, portcmd.NewPortLSCommand)
	do.Provide(i, appcmd.NewAppCommand)
	do.Provide(i, appcmd.NewAppLSCommand)
	do.Provide(i, appcmd.NewAppCreateCommand)
	do.Provide(i, entrypointcmd.NewEntryPointCommand)
	do.Provide(i, entrypointcmd.NewEntryPointLSCommand)
	do.Provide(i, entrypointcmd.NewEntryPointCreateCommand)
	do.Provide(i, infracmd.NewInfraCommand)
	do.Provide(i, infracmd.NewInfraLSCommand)
	do.Provide(i, infracmd.NewInfraCreateCommand)
	do.Provide(i, packagecmd.NewPackageCommand)
	do.Provide(i, packagecmd.NewPackageLSCommand)
	do.Provide(i, packagecmd.NewPackageCreateCommand)
	do.Provide(i, runnercmd.NewRunnerCommand)
	do.Provide(i, doctorcmd.NewDoctorCommand)
	do.Provide(i, treecmd.NewTreeCommand)

	do.Provide(i, projectservice.New)
	do.Provide(i, config.New)

	c, err := do.Invoke[*cli.CLI](i)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.Run(ctx)
	if err != nil {
		if !errors.Is(err, dto.ErrSuppressed) {
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

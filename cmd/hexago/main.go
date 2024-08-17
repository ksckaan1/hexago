package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ksckaan1/hexago/internal/domain/core/application/cli"
	projectservice "github.com/ksckaan1/hexago/internal/domain/core/service/project"
	"github.com/ksckaan1/hexago/internal/util"
	"github.com/samber/do"
)

func main() {
	ctx := context.Background()

	i := do.New()
	do.Provide(i, cli.NewCLI)
	do.Provide(i, cli.NewRootCommand)
	do.Provide(i, cli.NewInitCommand)
	do.Provide(i, projectservice.New)

	c, err := do.Invoke[*cli.CLI](i)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.Run(ctx)
	if err != nil {
		fmt.Println("")
		util.UILog(util.Error, util.UnwrapAllErrors(err).Error())
		fmt.Println("")
		os.Exit(1)
	}
}

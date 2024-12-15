package port

import "github.com/spf13/cobra"

type Commander interface {
	AddSubCommand(cmd Commander)
	Command() *cobra.Command
}

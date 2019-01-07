package cliexec

import (
	"github.com/spf13/cobra"
	"github.com/superioz/artemis/cliexec/status"
)

func Register(cmd *cobra.Command) {
	cmd.AddCommand(status.Cmd())
}

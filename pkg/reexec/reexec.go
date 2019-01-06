package reexec

import (
	"github.com/spf13/cobra"
)

func Register(cmd *cobra.Command) {
	cmd.AddCommand(statusCmd)
}

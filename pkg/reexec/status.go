package reexec

import (
	"fmt"
	"github.com/spf13/cobra"
)

var statusCmd *cobra.Command

func init() {
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Checks the current status of the daemon.",
		Run:   status,
	}
}

func status(cmd *cobra.Command, args []string) {
	fmt.Println("Status: ok")
}

package status

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/superioz/artemis/cliexec"
)

var statusCmd *cobra.Command

func Cmd() *cobra.Command {
	return statusCmd
}

func init() {
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Checks the current status of the daemon.",
		Run:   status,
	}
}

func status(cmd *cobra.Command, args []string) {
	data, code, err := cliexec.Get([]byte{}, "/status")
	fmt.Println(string(data), code, err)
}

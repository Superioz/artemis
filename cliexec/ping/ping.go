package ping

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/superioz/artemis/cliexec"
	"time"
)

var pingCmd *cobra.Command

func Cmd() *cobra.Command {
	return pingCmd
}

func init() {
	pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Displays the ping between cli and daemon.",
		Run:   ping,
	}
}

func ping(cmd *cobra.Command, args []string) {
	now := time.Now()
	_, code, err := cliexec.Get([]byte{}, "/")
	nowAfter := time.Now()

	if err != nil {
		fmt.Println("There is no daemon running.")
		return
	}
	if code != 200 {
		fmt.Println(fmt.Sprintf("Daemon returned unexpected code of %d", code))
		return
	}
	fmt.Println(fmt.Sprintf("Ping: +%s", nowAfter.Sub(now)))
}

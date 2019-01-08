package status

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/superioz/artemis/cliexec"
	"github.com/superioz/artemis/internal/dome"
	"strings"
	"time"
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

	if err != nil {
		fmt.Println("There is no daemon running.")
		return
	}
	if code != 200 {
		fmt.Println(fmt.Sprintf("Daemon returned unexpected code of %d", code))
		return
	}

	var status dome.Status
	err = json.Unmarshal(data, &status)
	if err != nil {
		fmt.Println("Daemon returned unexpected response.", err)
		return
	}

	var cStr string
	if status.BrokerConnected {
		cStr = "ok"
	} else {
		cStr = "nok"
	}

	now := time.Now()

	fmt.Println(fmt.Sprintf("Current status of artemis daemon %s: %s", status.Id.String(), cStr))
	fmt.Println(" ")
	fmt.Println(fmt.Sprintf("Running artemis version %s", status.Version))
	fmt.Println(fmt.Sprintf("Last activity: i +%s o +%s", now.Sub(status.LastReceived), now.Sub(status.LastSent)))
	fmt.Println(fmt.Sprintf("Connected with cluster of size %d", status.ClusterSize))
	fmt.Println(fmt.Sprintf("Running since: +%s", status.Runtime))
	fmt.Println(fmt.Sprintf("Ping: +%s", time.Now().Sub(status.Ping)))
	fmt.Println(" ")
	fmt.Println(fmt.Sprintf("Currently running as %s in TERM %d", strings.ToUpper(string(status.State)), status.Term))
	fmt.Println(fmt.Sprintf("Latest log view (10/%d):", len(status.Log)))

	fmt.Println(fmt.Sprintf("%-15s %-15s %-15s", "INDEX", "TERM", "COMMAND"))

	for _, e := range status.Log {
		fmt.Println(fmt.Sprintf("%-15d %-15d %-15s", e.Index, e.Term, string(e.Content)))
	}
}

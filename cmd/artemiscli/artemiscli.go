package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/superioz/artemis/appversion"
	"github.com/superioz/artemis/cliexec/ping"
	"github.com/superioz/artemis/cliexec/status"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "artemis",
		Short:   "Orchestration system based on raft.",
		Version: fmt.Sprintf("%s, build %s", appversion.Version, appversion.Build),
	}

	// commands
	rootCmd.AddCommand(status.Cmd())
	rootCmd.AddCommand(ping.Cmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

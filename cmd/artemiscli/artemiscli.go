package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/superioz/artemis/artemisversion"
	"github.com/superioz/artemis/pkg/reexec"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "artemis",
		Short:         "Orchestration system based on raft.",
		Version: fmt.Sprintf("%s, build %s", artemisversion.Version, artemisversion.Build),
	}
	reexec.Register(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

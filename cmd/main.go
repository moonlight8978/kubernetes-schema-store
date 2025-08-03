package cmd

import (
	"fmt"
	"os"

	"github.com/moonlight8978/kubernetes-schema-store/cmd/server"
	"github.com/moonlight8978/kubernetes-schema-store/cmd/version"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use: "kss",
	}

	serverCmd := server.NewCommand()
	versionCmd := version.NewCommand()

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package cmd

import (
	"os"

	"github.com/moonlight8978/kubernetes-schema-store/cmd/server"
	"github.com/moonlight8978/kubernetes-schema-store/cmd/sync"
	"github.com/moonlight8978/kubernetes-schema-store/cmd/version"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/log"
	"github.com/spf13/cobra"
)

func Execute() {
	// Initialize logging first
	log.InitializeDefault()

	rootCmd := &cobra.Command{
		Use:   "kss",
		Short: "Kubernetes Schema Store - Generate JSON schemas from Kubernetes OpenAPI specs",
		Long:  `A Kubernetes utilities application that generates object/CRDs JSON schema from OpenAPI documentation.`,
	}

	serverCmd := server.NewCommand()
	versionCmd := version.NewCommand()
	syncCmd := sync.NewCommand()

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

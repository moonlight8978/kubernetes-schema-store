package server

import (
	"fmt"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	serverConfig := config.ServerConfig{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Long:  `Start the server`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%+v", serverConfig)
		},
	}

	cmd.Flags().StringVarP(&serverConfig.AuthMethod, "auth-method", "m", "service-account", "			Authentication method, one of: service-account, kube-config")
	cmd.MarkFlagRequired("auth-method")

	cmd.Flags().StringVarP(&serverConfig.KubeConfig.Path, "kube-config", "p", "", "Kube Config Path")

	cmd.Flags().StringVarP(&serverConfig.Destination, "destination", "d", "", "Rclone destination")
	cmd.MarkFlagRequired("exporter")

	return cmd
}

package server

import (
	"github.com/moonlight8978/kubernetes-schema-store/operations"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/config"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	serverConfig := config.ServerConfig{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Long:  `Start the server`,
		Run: func(cmd *cobra.Command, args []string) {
			auth := kubernetes.Auth{
				Method: serverConfig.AuthMethod,
				KubeConfig: &config.KubeConfig{
					Path: serverConfig.KubeConfig.Path,
				},
			}

			cluster := auth.GetCluster()

			operations.Sync(cluster)
		},
	}

	cmd.Flags().StringVarP(&serverConfig.AuthMethod, "auth-method", "m", "service-account", "			Authentication method, one of: service-account, kubeconfig")
	cmd.MarkFlagRequired("auth-method")

	cmd.Flags().StringVarP(&serverConfig.KubeConfig.Path, "kubeconfig", "p", "", "Kube Config Path")

	cmd.Flags().StringVarP(&serverConfig.Destination, "destination", "d", "", "Rclone destination")
	cmd.MarkFlagRequired("exporter")

	return cmd
}

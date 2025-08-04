package sync

import (
	"github.com/moonlight8978/kubernetes-schema-store/operations"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/config"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	serverConfig := config.ServerConfig{}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync the schemas to storage",
		Long:  `Sync the schemas to storage`,
		Run: func(cmd *cobra.Command, args []string) {
			auth := kubernetes.Auth{
				Method: serverConfig.AuthMethod,
				KubeConfig: &config.KubeConfig{
					Path: serverConfig.KubeConfig.Path,
				},
			}

			cluster := auth.GetCluster()

			operations.Sync(cluster, serverConfig.Destination)
		},
	}

	cmd.Flags().StringVarP(&serverConfig.AuthMethod, "auth-method", "m", "in-cluster", "Authentication method, one of: in-cluster, kubeconfig")
	cmd.MarkFlagRequired("auth-method")

	cmd.Flags().StringVarP(&serverConfig.KubeConfig.Path, "kubeconfig", "p", "", "Kube Config Path")

	cmd.Flags().StringVarP(&serverConfig.Destination, "destination", "d", "", "Rclone destination")
	cmd.MarkFlagRequired("destination")

	return cmd
}

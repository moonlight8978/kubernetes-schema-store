package kubernetes

import (
	"fmt"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/config"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Auth struct {
	Method     string
	KubeConfig *config.KubeConfig
}

func (auth *Auth) GetCluster() *Cluster {
	restConfig, err := auth.BuildConfig()
	if err != nil {
		panic(err)
	}
	cluster := &Cluster{
		Config: restConfig,
	}
	cluster.Client, err = cluster.NewClient()
	if err != nil {
		panic(err)
	}
	cluster.ApiExtensionsClient, err = cluster.NewApiExtensionsClient()
	if err != nil {
		panic(err)
	}

	return cluster
}

func (auth *Auth) BuildConfig() (*rest.Config, error) {
	switch auth.Method {
	case "service-account":
		return rest.InClusterConfig()
	case "kubeconfig":
		return clientcmd.BuildConfigFromFlags("", auth.KubeConfig.Path)
	default:
		return nil, fmt.Errorf("invalid auth method: %s", auth.Method)
	}
}

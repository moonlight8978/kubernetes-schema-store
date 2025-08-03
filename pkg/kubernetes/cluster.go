package kubernetes

import (
	apiExtensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Cluster struct {
	Config              *rest.Config
	Client              *k8s.Clientset
	ApiExtensionsClient *apiExtensions.Clientset
}

func (cluster *Cluster) NewClient() (*k8s.Clientset, error) {
	return k8s.NewForConfig(cluster.Config)
}

func (cluster *Cluster) NewApiExtensionsClient() (*apiExtensions.Clientset, error) {
	return apiExtensions.NewForConfig(cluster.Config)
}

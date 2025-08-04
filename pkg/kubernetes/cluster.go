package kubernetes

import (
	"net/http"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/log"
	apiExtensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Cluster struct {
	Config              *rest.Config
	Client              *k8s.Clientset
	ApiExtensionsClient *apiExtensions.Clientset
	DynamicClient       *dynamic.DynamicClient
	DiscoveryClient     *discovery.DiscoveryClient
	HttpClient          *http.Client
}

func (cluster *Cluster) NewClient() (*k8s.Clientset, error) {
	client, err := k8s.NewForConfig(cluster.Config)
	if err != nil {
		log.Error("Failed to create Kubernetes client", "error", err)
		return nil, err
	}
	cluster.Client = client
	log.Debug("Created Kubernetes client successfully")
	return client, nil
}

func (cluster *Cluster) NewApiExtensionsClient() *apiExtensions.Clientset {
	client, err := apiExtensions.NewForConfig(cluster.Config)
	if err != nil {
		log.Error("Failed to create API extensions client", "error", err)
		panic(err)
	}
	cluster.ApiExtensionsClient = client
	log.Debug("Created API extensions client successfully")
	return client
}

func (cluster *Cluster) NewDiscoveryClient() *discovery.DiscoveryClient {
	client, err := discovery.NewDiscoveryClientForConfig(cluster.Config)
	if err != nil {
		log.Error("Failed to create discovery client", "error", err)
		panic(err)
	}
	cluster.DiscoveryClient = client
	log.Debug("Created discovery client successfully")
	return client
}

func (cluster *Cluster) NewDynamicClient() *dynamic.DynamicClient {
	client, err := dynamic.NewForConfig(cluster.Config)
	if err != nil {
		log.Error("Failed to create dynamic client", "error", err)
		panic(err)
	}
	cluster.DynamicClient = client
	log.Debug("Created dynamic client successfully")
	return client
}

func (cluster *Cluster) NewHttpClient() *http.Client {
	transport, err := rest.TransportFor(cluster.Config)
	if err != nil {
		log.Error("Failed to create HTTP client", "error", err)
		panic(err)
	}
	cluster.HttpClient = &http.Client{
		Transport: transport,
	}
	return cluster.HttpClient
}

// Deprecated: Use the individual New*Client methods that return errors instead
func ClientOrDie[T any](client *T, err error) *T {
	if err != nil {
		log.Fatal("Client creation failed", "error", err)
	}
	return client
}

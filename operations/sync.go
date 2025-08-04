package operations

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/fs"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/rclone"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/log"
	"github.com/rclone/rclone/fs/config/configfile"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/cel/openapi/resolver"
)

func Sync(cluster *kubernetes.Cluster, dst string) error {
	// Init rclone config file
	configfile.Install()

	// Init cluster
	cluster.NewDiscoveryClient()

	os.RemoveAll(fs.GetTmpDir())
	os.MkdirAll(fs.GetTmpDir(), 0755)

	err := ExportResources(cluster, dst)
	if err != nil {
		return err
	}

	os.RemoveAll(fs.GetTmpDir())

	return nil
}

func ExportResources(cluster *kubernetes.Cluster, dst string) error {
	apiResources, err := cluster.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		log.Error("Failed to get API resources", "error", err)
		return err
	}

	count := 0

	clientDiscoveryResolver := resolver.ClientDiscoveryResolver{
		Discovery: cluster.DiscoveryClient,
	}

	for _, apiResourceList := range apiResources {
		group, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			log.Error("Failed to parse group version", "error", err)
			continue
		}

		for _, apiResource := range apiResourceList.APIResources {
			groupVersionKind := kubernetes.ToGroupKindVersion(group, apiResource)
			metadata := kubernetes.ToSchemaMetadata(group, *groupVersionKind)

			groupSchema, err := clientDiscoveryResolver.ResolveSchema(*groupVersionKind)
			if err != nil {
				log.Error("Failed to resolve schema", "error", err)
				continue
			}

			jsonSchema, err := kubernetes.ToJson(groupSchema)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to convert schema to json %s %s %s", group.Group, group.Version, groupVersionKind.Kind), "error", err)
			}

			os.MkdirAll(filepath.Join(fs.GetTmpDir(), fs.GetSchemaDir(metadata)), 0755)
			os.WriteFile(filepath.Join(fs.GetTmpDir(), fs.GetSchemaPath(metadata)), jsonSchema, 0644)

			err = rclone.Sync(filepath.Join(fs.GetTmpDir(), fs.GetSchemaPath(metadata)), dst, fs.GetSchemaPath(metadata))
			if err != nil {
				log.Error("Failed to sync schema", "error", err)
			}
		}

		count += len(apiResourceList.APIResources)
	}

	log.Info(fmt.Sprintf("Exported %d resource schemas", count), slog.String("path", fs.GetTmpDir()))
	return nil
}

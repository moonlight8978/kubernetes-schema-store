package operations

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/fs"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/log"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/rclone"
	"github.com/rclone/rclone/fs/config/configfile"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Sync(cluster *kubernetes.Cluster, dst string) error {
	apiExtensionsClient := cluster.NewApiExtensionsClient()
	configfile.Install()
	// discoveryClient := cluster.NewDiscoveryClient()
	// dynamicClient := cluster.NewDynamicClient()

	err := ExportCRDs(apiExtensionsClient, dst)
	if err != nil {
		return err
	}

	return nil
}

func ExportCRDs(client *clientset.Clientset, dst string) error {
	crds, err := client.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Error("Failed to list CRDs", "error", err)
		return err
	}

	os.RemoveAll(fs.GetTmpDir())
	count := 0

	for _, crd := range crds.Items {
		schemas, err := kubernetes.ToJsonSchema(&crd)
		if err != nil {
			log.Error("Failed to convert CRD to JSON schema", "error", err)
			return err
		}

		for _, schema := range schemas {
			schemaJson, err := schema.ToJson()
			if err != nil {
				log.Error("Failed to convert JSON schema to JSON", "error", err)
				return err
			}

			os.MkdirAll(filepath.Join(fs.GetTmpDir(), fs.GetSchemaDir(schema)), 0755)
			os.WriteFile(filepath.Join(fs.GetTmpDir(), fs.GetSchemaPath(schema)), schemaJson, 0644)

			err = rclone.Sync(filepath.Join(fs.GetTmpDir(), fs.GetSchemaPath(schema)), dst, fs.GetSchemaPath(schema))

			if err != nil {
				log.Error("Failed to sync schema", "error", err)
			}
		}

		count += len(schemas)
	}

	log.Info(fmt.Sprintf("Exported schema %d", count), slog.String("path", fs.GetTmpDir()))

	return nil
}

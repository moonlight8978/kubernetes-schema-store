package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/fs"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes/schema"

	crdSchema "github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes/schema/crd"
	"github.com/moonlight8978/kubernetes-schema-store/pkg/log"
	"github.com/rclone/rclone/fs/config/configfile"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sSchema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kube-openapi/pkg/util/proto"
)

func Sync(cluster *kubernetes.Cluster, dst string) error {
	// Init rclone config file
	configfile.Install()

	// Init cluster
	cluster.NewApiExtensionsClient()
	cluster.NewDiscoveryClient()
	cluster.NewDynamicClient()
	cluster.NewHttpClient()

	os.RemoveAll(fs.GetTmpDir())

	err := ExportResources(cluster, dst)
	if err != nil {
		return err
	}

	return nil
}

func ExportCRDs(cluster *kubernetes.Cluster, dst string) error {
	crds, err := cluster.ApiExtensionsClient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Error("Failed to list CRDs", "error", err)
		return err
	}

	count := 0

	for _, crd := range crds.Items {
		schemas, err := crdSchema.ToJsonSchema(&crd)
		if err != nil {
			log.Error("Failed to convert CRD to JSON schema", "error", err)
			return err
		}

		for _, schema := range schemas {
			schemaJson, err := schema.ToJson()
			if err != nil {
				log.Error("Failed to convert JSON schema to JSON", "error", err)
			}

			os.MkdirAll(filepath.Join(fs.GetTmpDir(), fs.GetSchemaDir(schema.Metadata)), 0755)
			os.WriteFile(filepath.Join(fs.GetTmpDir(), fs.GetSchemaPath(schema.Metadata)), schemaJson, 0644)

			// err = rclone.Sync(filepath.Join(fs.GetTmpDir(), fs.GetSchemaPath(schema.Metadata)), dst, fs.GetSchemaPath(schema.Metadata))

			// if err != nil {
			// 	log.Error("Failed to sync schema", "error", err)
			// }
		}

		count += len(schemas)
	}

	log.Info(fmt.Sprintf("Exported schema %d", count), slog.String("path", fs.GetTmpDir()))

	return nil
}

func ExportResources(cluster *kubernetes.Cluster, dst string) error {
	// Get OpenAPI schema
	openapiSchema, err := cluster.DiscoveryClient.OpenAPISchema()
	if err != nil {
		log.Error("Failed to get OpenAPI schema", "error", err)
		return err
	}

	// Parse OpenAPI data
	models, err := proto.NewOpenAPIData(openapiSchema)
	if err != nil {
		log.Error("Failed to parse OpenAPI data", "error", err)
		return err
	}

	// Get all API resources
	apiResources, err := cluster.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		log.Error("Failed to get API resources", "error", err)
		return err
	}

	count := 0

	resp, err := cluster.NewHttpClient().Get(cluster.Config.Host + "/openapi/v3")
	if err != nil {
		log.Error("Failed to get OpenAPI schema", "error", err)
		return err
	}
	defer resp.Body.Close()

	var groupToOpenAPIServer map[string]map[string]struct {
		ServerRelativeURL string `json:"serverRelativeURL"`
	}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &groupToOpenAPIServer)

	os.MkdirAll(filepath.Join(fs.GetTmpDir()), 0755)
	// os.WriteFile(filepath.Join(fs.GetTmpDir(), "groupToOpenAPIServer.json"), body, 0644)

	for _, apiResource := range apiResources {
		group, err := k8sSchema.ParseGroupVersion(apiResource.GroupVersion)
		if err != nil {
			log.Error("Failed to parse group version", "error", err)
			continue
		}

		fmt.Printf("group: %s %s %s\n", group.Group, group.Version, apiResource.GroupVersion)

		var docUrl string
		var groupName string
		if group.Group == "" {
			docUrl = "api/v1"
			groupName = "core"
		} else {
			docUrl = "apis/" + apiResource.GroupVersion
			groupName = group.Group
		}

		openApiServer := groupToOpenAPIServer["paths"][docUrl]

		groupSchema, err := cluster.HttpClient.Get(cluster.Config.Host + openApiServer.ServerRelativeURL)
		if err != nil {
			log.Error("Failed to get group schema", "error", err)
			continue
		}
		defer groupSchema.Body.Close()

		groupSchemaBody, _ := io.ReadAll(groupSchema.Body)

		var openapiSchema struct {
			Components struct {
				Schemas map[string]map[string]any `json:"schemas"`
			} `json:"components"`
		}
		json.Unmarshal(groupSchemaBody, &openapiSchema)

		for schemaName, schemaDef := range openapiSchema.Components.Schemas {
			if schemaDef["x-kubernetes-group-version-kind"] == nil {
				continue
			}

			var matchedGroup bool

			for _, kubernetesGroupVersionKind := range schemaDef["x-kubernetes-group-version-kind"].([]any) {
				kubernetesGroupVersionKindMap := kubernetesGroupVersionKind.(map[string]interface{})
				if kubernetesGroupVersionKindMap["group"] == group.Group && kubernetesGroupVersionKindMap["version"] == group.Version {
					matchedGroup = true
					break
				}
			}

			if !matchedGroup {
				continue
			}

			model := models.LookupModel(schemaName)
			if model == nil {
				continue
			}

			schemaJson, err := json.MarshalIndent(schemaDef, "", "  ")

			if err != nil {
				log.Error("Failed to marshal schema", "error", err)
				continue
			}
			nameParts := strings.Split(schemaName, ".")
			metadata := &schema.SchemaMetadata{
				Package: groupName,
				Version: group.Version,
				Name:    strings.ToLower(nameParts[len(nameParts)-1]),
			}

			os.MkdirAll(filepath.Join(fs.GetTmpDir(), fs.GetSchemaDir(metadata)), 0755)
			os.WriteFile(filepath.Join(fs.GetTmpDir(), fs.GetSchemaPath(metadata)), schemaJson, 0644)
		}
	}

	log.Info(fmt.Sprintf("Exported %d resource schemas", count), slog.String("path", fs.GetTmpDir()))
	return nil
}

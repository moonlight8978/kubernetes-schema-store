package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Sync(cluster *kubernetes.Cluster) error {
	apiExtensionsClient := cluster.ApiExtensionsClient
	crds, err := apiExtensionsClient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		for _, crd := range crds.Items {
			for _, version := range crd.Spec.Versions {
				openApiSchema := version.Schema.OpenAPIV3Schema
				if openApiSchema == nil {
					return fmt.Errorf("openapi schema for %s is empty", crd.Name)
				}

				schema, err := json.MarshalIndent(openApiSchema, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal OpenAPI schema for %s: %v", crd.Name, err)
				}

				parts := strings.Split(crd.Name, ".")
				name := parts[0]
				org := strings.Join(parts[1:], ".")

				fmt.Printf("%s %s %s\n", org, name, version.Name)

				os.MkdirAll(filepath.Join("tmp", org), 0755)
				os.WriteFile(filepath.Join("tmp", org, strings.Join([]string{name, version.Name}, "-")+".json"), schema, 0644)
			}
		}
	}

	return nil
}

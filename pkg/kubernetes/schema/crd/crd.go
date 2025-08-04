package crd

import (
	"encoding/json"
	"fmt"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes/schema"
)

type CRDJsonSchema struct {
	Metadata *schema.SchemaMetadata
	Schema   *apiextensionsv1.JSONSchemaProps
}

func ToJsonSchema(crd *apiextensionsv1.CustomResourceDefinition) ([]*CRDJsonSchema, error) {
	var schemas []*CRDJsonSchema

	for _, version := range crd.Spec.Versions {
		openApiSchema := version.Schema.OpenAPIV3Schema
		if openApiSchema == nil {
			return nil, fmt.Errorf("openapi schema for %s is empty", crd.Name)
		}

		parts := strings.Split(crd.Name, ".")
		name := parts[0]
		org := strings.Join(parts[1:], ".")

		schemas = append(schemas, &CRDJsonSchema{
			Metadata: &schema.SchemaMetadata{
				Package: org,
				Version: version.Name,
				Name:    name,
			},
			Schema: openApiSchema,
		})
	}

	return schemas, nil
}

func (js *CRDJsonSchema) ToJson() ([]byte, error) {
	return json.MarshalIndent(js.Schema, "", "  ")
}

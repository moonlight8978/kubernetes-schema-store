package resource

import (
	"encoding/json"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes/schema"
)

type ResourceJsonSchema struct {
	Metadata *schema.SchemaMetadata
	Schema   *v1.APIResource
}

func ToJsonSchema(group *v1.APIResourceList, resource *v1.APIResource) (*ResourceJsonSchema, error) {
	parts := strings.Split(group.GroupVersion, "/")
	org := parts[0]
	version := parts[1]
	name := strings.ToLower(resource.Kind)

	return &ResourceJsonSchema{
		Metadata: &schema.SchemaMetadata{
			Package: org,
			Version: version,
			Name:    name,
		},
		Schema: resource,
	}, nil
}

func (js *ResourceJsonSchema) ToJson() ([]byte, error) {
	return json.MarshalIndent(js.Schema, "", "  ")
}

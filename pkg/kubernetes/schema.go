package kubernetes

import (
	"encoding/json"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

func ToSchemaMetadata(gv schema.GroupVersion, gvk schema.GroupVersionKind) *SchemaMetadata {
	groupName := gv.Group
	if len(gv.Group) == 0 {
		groupName = "core"
	}
	return &SchemaMetadata{
		Package: groupName,
		Version: gv.Version,
		Name:    gvk.Kind,
	}
}

func ToGroupKindVersion(gv schema.GroupVersion, res v1.APIResource) *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   gv.Group,
		Version: gv.Version,
		Kind:    strings.ToLower(res.Kind),
	}
}

func ToJson(schema *spec.Schema) ([]byte, error) {
	return json.Marshal(schema)
}

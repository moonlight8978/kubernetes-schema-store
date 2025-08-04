package fs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/moonlight8978/kubernetes-schema-store/pkg/kubernetes"
)

func GetTmpDir() string {
	return filepath.Join(os.TempDir(), "kubernetes-schema-store")
}

func GetSchemaPath(schema *kubernetes.SchemaMetadata) string {
	return filepath.Join(GetSchemaDir(schema), strings.Join([]string{schema.Name, schema.Version}, "-")+".json")
}

func GetSchemaDir(schema *kubernetes.SchemaMetadata) string {
	return filepath.Join(schema.Package)
}

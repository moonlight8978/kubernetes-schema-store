package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	date    = "unknown"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kubernetes-schema-store %s build date %s\n", version, date)
		},
	}
}

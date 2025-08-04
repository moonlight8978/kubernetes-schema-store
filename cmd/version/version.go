package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Author  = "moonlight8978"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kubernetes-schema-store %s by %s\n", Version, Author)
		},
	}
}

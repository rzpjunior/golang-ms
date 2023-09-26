package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Edenlabs",
	Long:  `All software has versions. This is Edenlabs's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Edenlabs Library & Tools For Supporting Microservices v.1.0.0")
	},
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help about any command",
	Long:  `Help about any command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Edenlabs Library & Tools For Supporting Microservices v.1.0.0")
	},
}

func help(*cobra.Command, []string) {

}

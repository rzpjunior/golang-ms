package cmd

import (
	"os"

	edenlabs "git.edenfarm.id/edenlabs/edenlabs/server"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
			examples and usage of using your application.`,
	Run: func(_ *cobra.Command, _ []string) {
		start()
	},
}

//to reload
// main creating new instances application and serving application server.
func start() {
	var err error
	if global.Setup.Common, err = edenlabs.Start(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
		return
	}
	server.StartRestServer()
}

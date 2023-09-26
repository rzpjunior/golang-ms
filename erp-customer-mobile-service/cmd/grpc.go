package cmd

import (
	"os"

	edenlabs "git.edenfarm.id/edenlabs/edenlabs/server"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
			examples and usage of using your application.`,
	Run: func(_ *cobra.Command, _ []string) {
		startGrpc()
	},
}

func startGrpc() {
	var err error
	if global.Setup.Common, err = edenlabs.Start(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
		return
	}

	server.StartGrpcServer()
}

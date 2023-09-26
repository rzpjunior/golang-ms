package cmd

import (
	"os"

	edenlabs "git.edenfarm.id/edenlabs/edenlabs/server"
	grpcProvider "git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/provider"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/server"
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

	global.Setup.Common.Client.AuditServiceGrpc, err = grpcProvider.ConnectAuditServiceGrpc(grpcProvider.ProviderOptions{
		Common: global.Setup.Common,
	})
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
		return
	}

	global.Setup.Common.Client.ConfigurationServiceGrpc, err = grpcProvider.ConnectConfigurationServiceGrpc(grpcProvider.ProviderOptions{
		Common: global.Setup.Common,
	})
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
		return
	}

	global.Setup.Common.Client.BridgeServiceGrpc, err = grpcProvider.ConnectBridgeServiceGrpc(grpcProvider.ProviderOptions{
		Common: global.Setup.Common,
	})
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
		return
	}

	server.StartGrpcServer()
}

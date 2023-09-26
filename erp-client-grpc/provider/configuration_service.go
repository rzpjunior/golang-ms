package provider

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/client"
)

func ConnectConfigurationServiceGrpc(opt ProviderOptions) (configurationClientService client.IConfigurationServiceGrpc, err error) {
	// setup client configuration service
	configurationClientService, err = client.NewConfigurationServiceGrpc(client.ConfigurationServiceGrpcOption{
		Host:                  opt.Common.Env.GetString("client.configuration_service_grpc.host"),
		Port:                  opt.Common.Env.GetInt("client.configuration_service_grpc.port"),
		Timeout:               opt.Common.Env.GetDuration("client.configuration_service_grpc.timeout"),
		MaxConcurrentRequests: opt.Common.Env.GetInt("client.configuration_service_grpc.max_concurrent_requests"),
		ErrorPercentThreshold: opt.Common.Env.GetInt("client.configuration_service_grpc.error_percent_threshold"),
		Tls:                   opt.Common.Env.GetBool("client.configuration_service_grpc.tls"),
		PemPath:               opt.Common.Env.GetString("client.configuration_service_grpc.pem_tls_path"),
		Secret:                opt.Common.Env.GetString("client.configuration_service_grpc.secret"),
		Realtime:              opt.Common.Env.GetBool("client.configuration_service_grpc.realtime"),
	})
	if err != nil {
		opt.Common.Logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to connect client, error connect to configuration service | %v", err)).Print()
	}
	return
}

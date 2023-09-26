package provider

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/client"
)

func ConnectBridgeServiceGrpc(opt ProviderOptions) (bridgeClientService client.IBridgeServiceGrpc, err error) {
	// setup client bridge service
	bridgeClientService, err = client.NewBridgeServiceGrpc(client.BridgeServiceGrpcOption{
		Host:                  opt.Common.Env.GetString("client.bridge_service_grpc.host"),
		Port:                  opt.Common.Env.GetInt("client.bridge_service_grpc.port"),
		Timeout:               opt.Common.Env.GetDuration("client.bridge_service_grpc.timeout"),
		MaxConcurrentRequests: opt.Common.Env.GetInt("client.bridge_service_grpc.max_concurrent_requests"),
		ErrorPercentThreshold: opt.Common.Env.GetInt("client.bridge_service_grpc.error_percent_threshold"),
		Tls:                   opt.Common.Env.GetBool("client.bridge_service_grpc.tls"),
		PemPath:               opt.Common.Env.GetString("client.bridge_service_grpc.pem_tls_path"),
		Secret:                opt.Common.Env.GetString("client.bridge_service_grpc.secret"),
		Realtime:              opt.Common.Env.GetBool("client.bridge_service_grpc.realtime"),
	})
	if err != nil {
		opt.Common.Logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to connect client, error connect to bridge service | %v", err)).Print()
	}
	return
}

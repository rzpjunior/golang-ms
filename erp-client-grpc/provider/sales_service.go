package provider

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/client"
)

func ConnectSalesServiceGrpc(opt ProviderOptions) (SalesClientService client.ISalesServiceGrpc, err error) {
	// setup client Sales service
	SalesClientService, err = client.NewSalesServiceGrpc(client.SalesServiceGrpcOption{
		Host:                  opt.Common.Env.GetString("client.Sales_service_grpc.host"),
		Port:                  opt.Common.Env.GetInt("client.Sales_service_grpc.port"),
		Timeout:               opt.Common.Env.GetDuration("client.Sales_service_grpc.timeout"),
		MaxConcurrentRequests: opt.Common.Env.GetInt("client.Sales_service_grpc.max_concurrent_requests"),
		ErrorPercentThreshold: opt.Common.Env.GetInt("client.Sales_service_grpc.error_percent_threshold"),
		Tls:                   opt.Common.Env.GetBool("client.Sales_service_grpc.tls"),
		PemPath:               opt.Common.Env.GetString("client.Sales_service_grpc.pem_tls_path"),
		Secret:                opt.Common.Env.GetString("client.Sales_service_grpc.secret"),
		Realtime:              opt.Common.Env.GetBool("client.Sales_service_grpc.realtime"),
	})
	if err != nil {
		opt.Common.Logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to connect client, error connect to Sales service | %v", err)).Print()
	}
	return
}
